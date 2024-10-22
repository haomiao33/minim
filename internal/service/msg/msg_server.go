package msg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/hashicorp/consul/api"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"im/internal/command"
	"im/internal/dao"
	"im/internal/db"
	"im/internal/service/msg/config"
	"im/pb"
	"im/pkg/kafkaclient"
	"im/pkg/redisclient"
	"log"
	"net"
	"time"
)

type MsgServer struct {
	pb.UnimplementedMsgServiceServer
	kafkaClient *kafkaclient.KafkaProductClient
	redisClient *redisclient.RedisClient
	ctx         context.Context
}

func NewMsgServer(ctx context.Context) *MsgServer {
	redisClient := redisclient.NewRedisClient(ctx,
		config.Config.Redis.Addr,
		config.Config.Redis.Password)

	return &MsgServer{
		kafkaClient: kafkaclient.NewKafkaProductClient(
			ctx,
			config.Config.Kafka.Addresses),
		redisClient: redisClient,
		ctx:         ctx,
	}
}

// 单聊
func (s *MsgServer) handleSingleMsg(ctx context.Context, msg *pb.ImMsgRequest) (*pb.ImMsgResponse, error) {
	getMsg, err := dao.MsgDao.GetMsg(db.Db, msg.MsgId)
	if err != nil {
		logging.Errorf("--- get msg failed:%s ---", msg.MsgId)
		return nil, err
	}
	if getMsg != nil {
		logging.Infof("--- msg already exist:%s ---", msg.MsgId)
		return nil, errors.New("msg already exist")
	}

	bigId := msg.FromId
	if msg.ToId > bigId {
		bigId = msg.ToId
	}
	smallId := msg.ToId
	if smallId > msg.FromId {
		smallId = msg.FromId
	}

	//lock
	key := fmt.Sprintf("lock:msg:%d:%d:%d", msg.MsgType, bigId, smallId)
	mutex := s.redisClient.RedisLock.NewMutex(key, redsync.WithExpiry(5*time.Second))
	if err := mutex.Lock(); err != nil {
		logging.Errorf("--- lock failed:%s ---", msg.MsgId)
	}

	//事务开启
	tx := db.Db.Begin()

	defer func() {
		if err != nil {
			logging.Errorf("--- rollback:%s ---", msg.MsgId)
			tx.Rollback()
		}
		if ok, err := mutex.Unlock(); !ok || err != nil {
			logging.Errorf("--- unlock failed:%s ---", msg.MsgId)
		}
	}()

	//查询会话关系
	conversation, err := dao.ConversationDao.GetConversation(tx, msg.ChatType, bigId, smallId)
	if err != nil {
		logging.Errorf("--- get conversation failed:%s ---", msg.MsgId)
		return nil, err
	}
	if conversation == nil {
		addConversation, err := dao.ConversationDao.AddConversation(tx, msg.ChatType, bigId, smallId)
		if err != nil {
			logging.Errorf("--- create conversation failed:%s ---", msg.MsgId)
			return nil, err
		}
		conversation = addConversation
	}

	//获取序号
	seqKey := fmt.Sprintf("sequence:%d", conversation.ID)
	sequence, err := s.redisClient.GetSequence(seqKey, conversation.Sequence)
	if err != nil {
		logging.Errorf("--- get sequence failed:%s ---", msg.MsgId)
		return nil, err
	}
	//添加消息
	addMsg, err := dao.MsgDao.AddMsg(tx, conversation.ID, sequence, msg.MsgId,
		int(msg.ChatType), int(msg.MsgType),
		msg.FromId, msg.ToId, msg.Message, msg.Ts, 0)
	if err != nil {
		logging.Errorf("--- add msg failed:%s ---", msg.MsgId)
		return nil, err
	}

	//更新会话序号
	err = dao.ConversationDao.UpdateConversation(tx, conversation.ID, sequence)
	if err != nil {
		logging.Errorf("--- update conversation failed:%s ---", msg.MsgId)
		return nil, err
	}

	//发送者session
	senderSession, err := dao.RecentSessionDao.Get(tx, msg.ChatType, msg.FromId, msg.ToId)
	if err != nil {
		logging.Errorf("--- get recent session failed:%s ---", msg.MsgId)
		return nil, err
	}
	if senderSession == nil {
		err := dao.RecentSessionDao.Add(tx,
			msg.ChatType, msg.FromId, msg.ToId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logging.Errorf("--- add recent session failed:%s ---", msg.MsgId)
			return nil, err
		}
	} else {
		err := dao.RecentSessionDao.Update(tx,
			msg.ChatType, msg.FromId, msg.ToId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logging.Errorf("--- update recent session failed:%s ---", msg.MsgId)
			return nil, err
		}
	}

	//添加接收者session
	receiverSession, err := dao.RecentSessionDao.Get(tx, msg.ChatType, msg.ToId, msg.FromId)
	if err != nil {
		logging.Errorf("--- get recent session failed:%s ---", msg.MsgId)
		return nil, err
	}

	if receiverSession == nil {
		err := dao.RecentSessionDao.Add(tx,
			msg.ChatType, msg.ToId, msg.FromId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logging.Errorf("--- add recent session failed:%s ---", msg.MsgId)
			return nil, err
		}
	} else {
		err := dao.RecentSessionDao.Update(tx,
			msg.ChatType, msg.ToId, msg.FromId,
			addMsg.ID, addMsg.Content, time.UnixMilli(msg.Ts))
		if err != nil {
			logging.Errorf("--- update recent session failed:%s ---", msg.MsgId)
			return nil, err
		}
	}

	logging.Infof("send msg success  msgId:%s,chatType:%d, msgType:%d, sequence:%d",
		msg.MsgId, msg.ChatType, msg.MsgType, sequence)
	//提交事务
	tx.Commit()

	//发送消息到kafka，进行后续推送
	//相同会话发送到同一个分区
	partition := int32(conversation.ID % int64(config.Config.Kafka.MsgPartitionCount))

	marshal, _ := json.Marshal(addMsg)
	err = s.kafkaClient.ProductMessage(
		config.Config.Kafka.MsgTopic,
		partition,
		marshal)
	if err != nil {
		logging.Errorf("--- send message to kafka error: %v， msgId:%s ---", err, msg.MsgId)
		//这里后面有其他兜底（客户端主动pull），这里忽略
		err = nil
	} else {
		logging.Infof("add to push mq  success:%s", msg.MsgId)
	}

	return &pb.ImMsgResponse{
		Type: command.COMMAND_RESP_TYPE_MSG_ACK,
		Code: command.COMMAND_RESP_CODE_SUCESS,
		Msg:  "success",
		Data: &pb.ImMsgResponseData{
			MsgId:          msg.MsgId,
			Sequence:       sequence,
			ConversationId: conversation.ID,
		},
	}, nil
}

func (s *MsgServer) SendMessage(ctx context.Context, msg *pb.ImMsgRequest) (*pb.ImMsgResponse, error) {
	logging.Infof("MsgServer on recv msg : %v", msg)

	resp, err := s.handleSingleMsg(ctx, msg)
	if err != nil {
		logging.Errorf("send message to kafka error: %v", err)
		return &pb.ImMsgResponse{
			Type: command.COMMAND_RESP_TYPE_MSG_ACK,
			Code: command.COMMAND_RESP_CODE_SERVER_ERR,
			Msg:  "error to send to kafka",
			Data: &pb.ImMsgResponseData{
				MsgId: msg.MsgId,
			},
		}, nil
	}
	return resp, nil
}

func (s *MsgServer) Run() {
	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d",
		cfg.ListenHost,
		cfg.ListenPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	logging.Infof("starting gRPC server on %s", addr)
	grpcServer := grpc.NewServer()
	pb.RegisterMsgServiceServer(grpcServer, s)

	// 创建 Consul 客户端
	consulCfg := api.DefaultConfig()
	consulCfg.Address = config.Config.Consul.Address
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// 创建服务注册信息
	//randomId, _ := uuid.NewUUID()
	registration := &api.AgentServiceRegistration{
		ID:      "MsgService-1",
		Name:    "MsgService",
		Address: config.Config.Rpc.ListenHost,
		Port:    config.Config.Rpc.ListenPort,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", config.Config.Rpc.ListenHost, config.Config.Rpc.ListenPort),
			GRPCUseTLS:                     false,
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}
	// 注册服务到 Consul
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	log.Println("gRPC server is running on: ", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
