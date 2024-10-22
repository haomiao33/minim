package msg_push

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	_ "google.golang.org/grpc/health"
	"im/internal/command"
	"im/internal/model"
	"im/internal/service/msg_push/config"
	"im/pb"
	"im/pkg/kafkaclient"
	"im/pkg/rpcclient"
	"log"
	"sync"
	"time"
)

type MsgPushServer struct {
	kafkaClient *kafkaclient.KafkaConsumerClient
	onlineRpc   pb.OnlineServiceClient
	ctx         context.Context

	consulClient *api.Client

	//所有的登录服务客户端
	loginServiceMap sync.Map
}

func NewMsgPushServer(ctx context.Context) *MsgPushServer {
	// 创建 Consul 客户端
	consulCfg := api.DefaultConfig()
	consulCfg.Address = config.Config.Consul.Address
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	return &MsgPushServer{
		ctx:          ctx,
		consulClient: consulClient,
		kafkaClient: kafkaclient.NewKafkaConsumerClient(
			ctx,
			config.Config.Kafka.Addresses,
			config.Config.Kafka.MsgTopicGroup),
		onlineRpc: rpcclient.NewOnlineRpcClient(
			ctx,
			config.Config.Consul.Address,
			"OnlineService",
		),
	}
}

// 获取所有登录服务实例地址
func (s *MsgPushServer) getLoginServiceAddresses(serviceName string, serviceId string) (string, error) {
	instances, _, err := s.consulClient.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to get service instances: %v", err)
	}

	for _, instance := range instances {
		if instance.ServiceID == serviceId {
			return fmt.Sprintf("%s:%d", instance.ServiceAddress,
				instance.ServicePort), nil
		}
	}
	return "", fmt.Errorf("service instance not found")
}

// 推送消息到在线用户
func (s *MsgPushServer) pushMessageToOnlineUser(msg model.ImMsg, serviceId string) error {
	var loginClient pb.LoginServiceClient
	// 检查是否已有连接
	value, ok := s.loginServiceMap.Load(serviceId)
	if !ok {
		// 没有连接，获取服务地址
		address, err := s.getLoginServiceAddresses("LoginService", serviceId) // 替换为你的服务名称
		if err != nil {
			logging.Errorf("failed to get login service addresses: %s", err.Error())
			return err
		}
		///这里用直连的grpc 客户端，不用consul
		loginClient = rpcclient.NewLoginDefaultRpcClient(s.ctx, address, "LoginService")

		// 存储连接到 loginService
		s.loginServiceMap.Store(serviceId, loginClient)
	} else {
		// 已有连接，直接推送
		loginClient = value.(pb.LoginServiceClient)
	}

	// 发送消息
	cmd := command.ImMsgCommandResp{
		Type: command.COMMAND_TYPE_MSG_SYNC_NOTIFY,
		Code: command.COMMAND_RESP_CODE_SUCESS,
		Msg:  "msg",
		Data: command.ImMsgSyncNotifyCommand{
			FromId:         msg.FromID,
			ToId:           msg.ToID,
			MsgId:          msg.ID,
			ChatType:       msg.ChatType,
			MsgType:        msg.MsgType,
			ConversationId: msg.ConversationID,
			Sequence:       msg.Sequence,
		},
	}
	marshal, _ := json.Marshal(cmd)
	_, err := loginClient.PushMsg(s.ctx, &pb.PushRequest{
		UserId: msg.ToID,
		Data:   string(marshal),
	})
	if err != nil {
		logging.Errorf("failed to push online message:%s to %s: %v",
			msg.ID, msg.ToID, err)
		return err
	}
	return nil
}

func (s *MsgPushServer) OnPushMsg(ctx context.Context, data []byte) error {
	var msg model.ImMsg
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logging.Errorf("--- json unmarshal failed:%s, data:%s ---", err, string(data))
		return err
	}
	logging.Infof("start push msg : msgId:%s, chatType:%d, "+
		"msgType:%d, conversationId:%d, sequence:%d",
		msg.ID, msg.ChatType, msg.MsgType, msg.ConversationID, msg.Sequence)

	resp, err := s.onlineRpc.GetOnlineUser(ctx, &pb.GetOnlineUserRequest{UserId: msg.ToID})
	if err != nil {
		logging.Infof("user is outline,  toId:%d, msgId:%s", msg.ToID, msg.ID)
		//todo 离线推送

		return nil
	}

	logging.Infof("user:%d is online, start push msg:%s", msg.ToID, msg.ID)
	err = s.pushMessageToOnlineUser(msg, resp.ServerId)
	if err != nil {
		logging.Infof("user:%d is outline, msgId:%s", msg.ToID, msg.ID)
		//todo 离线推送

		return nil
	}
	logging.Infof("end push online message:%s to %d",
		msg.ID, msg.ToID)

	return nil
}

func (s *MsgPushServer) Run() error {
	// 监听服务离线
	go func() {
		// 开始监视
		for {
			services, _, err := s.consulClient.Health().Service("LoginService", "",
				true, &api.QueryOptions{
					WaitTime: 12 * time.Second,
				})
			if err != nil {
				//清空
				logging.Errorf("consul service connect error")
				time.Sleep(2 * time.Second)
				continue
			}

			s.loginServiceMap.Range(func(key, value interface{}) bool {
				for _, service := range services {
					if service.Service.ID == key.(string) {
						return true
					}
				}
				// 移除已下线的服务
				s.loginServiceMap.Delete(key)
				logging.Infof("service instance %s is offline", key.(string))
				return true
			})
			time.Sleep(5 * time.Second) // 适当的重试间隔
		}
	}()

	return s.kafkaClient.StartConsume([]string{config.Config.Kafka.MsgTopic},
		func(data []byte) error {
			return s.OnPushMsg(s.ctx, data)
		})
}
