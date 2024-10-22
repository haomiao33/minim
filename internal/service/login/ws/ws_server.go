package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"im/internal/command"
	"im/internal/service/login"
	"im/internal/service/login/config"
	"im/pb"
	"im/pkg/rpcclient"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 心跳超时时间 毫秒
const HEART_TIME_OUT_MILLI_SECOND = 21 * 1000

type OnlineUser struct {
	UserId int64
	Conn   gnet.Conn
	Ts     int64
}

type WsServer struct {
	uniqueId string
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
	connected int64
	users     sync.Map

	ctx        context.Context
	cancelFunc context.CancelFunc

	//router rpc client
	//routerClient rpcclient.IRouterRpcClient

	msgRpcClient    pb.MsgServiceClient
	onlineRpcClient pb.OnlineServiceClient

	//给外部提供一个推送给客户端的推送接口
	pb.LoginServiceServer
}

func NewWsServer(ctx context.Context) login.ServerBase {
	wsConfig := config.Config.Websocket
	wsCtx, cancelFunc := context.WithCancel(ctx)

	newUUID, _ := uuid.NewUUID()

	return &WsServer{
		uniqueId: newUUID.String(),
		addr: fmt.Sprintf(
			"tcp4://%s:%d",
			wsConfig.ListenHost,
			wsConfig.ListenPort),
		multicore:  wsConfig.Multicore,
		ctx:        wsCtx,
		cancelFunc: cancelFunc,
		msgRpcClient: rpcclient.NewMsgRpcClient(wsCtx,
			config.Config.Consul.Address, "MsgService"),
		onlineRpcClient: rpcclient.NewOnlineRpcClient(wsCtx,
			config.Config.Consul.Address, "OnlineService"),
	}
}

func (ws *WsServer) Run() {
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
	pb.RegisterLoginServiceServer(grpcServer, ws)

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
		ID:      fmt.Sprintf("LoginService-%s", ws.uniqueId),
		Name:    "LoginService",
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

	go func() {
		log.Println("gRPC server is running on: ", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Fatal(gnet.Run(ws, ws.addr,
		gnet.WithMulticore(ws.multicore),
		gnet.WithReusePort(true),
		gnet.WithTicker(true)))
}

func (ws *WsServer) OnBoot(eng gnet.Engine) gnet.Action {
	ws.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n",
		ws.multicore, ws.addr)

	go func() {
		for {
			select {
			case <-ws.ctx.Done():
				logging.Infof("im heart checker exit!!")
				return
			default:
				time.Sleep(10 * time.Second)

				now := time.Now().UnixMilli()
				ws.users.Range(func(key, value interface{}) bool {
					onlineUser := value.(*OnlineUser)
					if now-onlineUser.Ts > HEART_TIME_OUT_MILLI_SECOND {
						logging.Infof("user[%v] heartbeat timeout ts:%v",
							onlineUser.UserId, now-onlineUser.Ts)
						onlineUser.Conn.Close()
						ws.users.Delete(key)

						//下线
						ws.onlineRpcClient.OutlineUser(ws.ctx, &pb.GetOnlineUserRequest{
							UserId: onlineUser.UserId,
						})
					}
					return true
				})
			}
		}
	}()

	return gnet.None
}

func (wss *WsServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(wsCodec))
	atomic.AddInt64(&wss.connected, 1)
	logging.Infof("conn[%v] connected", c.RemoteAddr().String())
	return nil, gnet.None
}

func (wss *WsServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	atomic.AddInt64(&wss.connected, -1)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())

	//查询用户
	user, ok := wss.users.Load(c.RemoteAddr().String())
	if ok {
		//下线
		wss.onlineRpcClient.OutlineUser(wss.ctx, &pb.GetOnlineUserRequest{
			UserId: user.(*OnlineUser).UserId,
		})
	}

	return gnet.None
}

func (wss *WsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*wsCodec)
	if ws.readBufferBytes(c) == gnet.Close {
		return gnet.Close
	}
	ok, action := ws.upgrade(c)
	if !ok {
		return
	}

	if ws.buf.Len() <= 0 {
		return gnet.None
	}
	messages, err := ws.Decode(c)
	if err != nil {
		return gnet.Close
	}
	if messages == nil {
		return
	}
	for _, message := range messages {
		payloadStr := string(message.Payload)
		var cmd command.ImCommand
		err := json.Unmarshal(message.Payload, &cmd)
		if err != nil {
			logging.Errorf("conn[%v] [err=%v] [params=%s]",
				c.RemoteAddr().String(), err.Error(), payloadStr)
			return gnet.Close
		}

		switch cmd.Type {
		case command.COMMAND_TYPE_HEARTBEAT:
			user, ok := wss.users.Load(c.RemoteAddr().String())
			if ok {
				//更新下线时间
				onlineUser := user.(*OnlineUser)
				onlineUser.Ts = time.Now().UnixMilli()
				//心跳
				wss.onlineRpcClient.UpdateOnlineUser(wss.ctx, &pb.UpdateOnlineUserRequest{
					UserId:       onlineUser.UserId,
					ServerId:     fmt.Sprintf("LoginService-%s", wss.uniqueId),
					LastUpdateTs: onlineUser.Ts,
				})
			}
			return gnet.None
		case command.COMMAND_TYPE_LOGIN_REQ:
			logging.Infof("COMMAND_TYPE_LOGIN:%s", payloadStr)
			//用户登陆
			body, _ := json.Marshal(cmd.Data)
			return wss.OnLogin(c, body)
		case command.COMMAND_TYPE_MSG:
			logging.Infof("COMMAND_TYPE_MSG :%s", payloadStr)
			body, _ := json.Marshal(cmd.Data)
			return wss.OnImMsg(c, body)
		case command.COMMAND_TYPE_MSG_SYNC:
			logging.Infof("COMMAND_TYPE_MSG_SYNC :%s", payloadStr)
			body, _ := json.Marshal(cmd.Data)
			return wss.OnMsgSync(c, body)
		default:
			logging.Errorf("conn[%v] [err=%v] [params=%s]",
				c.RemoteAddr().String(), "unknown command", payloadStr)
			return gnet.Close
		}

		//// This is the echo server
		//err = wsutil.WriteServerMessage(c, message.OpCode, message.Payload)
		//if err != nil {
		//	logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		//	return gnet.Close
		//}
	}
	return gnet.None
}

// 登陆
func (wss *WsServer) OnLogin(c gnet.Conn, body []byte) gnet.Action {
	var loginCmd command.ImLoginCommandReq
	err := json.Unmarshal(body, &loginCmd)
	if err != nil {
		logging.Errorf("command.ImLoginCommandReq parse err,%v", err)
		return gnet.Close
	}

	//登陆存储下
	wss.users.Store(c.RemoteAddr().String(), &OnlineUser{
		UserId: loginCmd.UserId,
		Conn:   c,
		Ts:     time.Now().UnixMilli(),
	})

	//登陆成功
	cmd := command.ImLoginCommandResp{
		Type: command.COMMAND_TYPE_LOGIN_RESP,
		Code: 200,
		Msg:  "success",
	}
	marshal, _ := json.Marshal(cmd)
	wsutil.WriteServerMessage(c, ws.OpText, marshal)

	//更新状态
	wss.onlineRpcClient.UpdateOnlineUser(wss.ctx, &pb.UpdateOnlineUserRequest{
		UserId:       loginCmd.UserId,
		ServerId:     fmt.Sprintf("LoginService-%s", wss.uniqueId),
		LastUpdateTs: time.Now().UnixMilli(),
	})
	return gnet.None
}

// 私聊消息
func (wss *WsServer) OnImMsg(c gnet.Conn, body []byte) gnet.Action {
	var msgCmd command.ImMsgCommandReq
	err := json.Unmarshal(body, &msgCmd)
	if err != nil {
		logging.Errorf("bad ImMsgCommandReq, %v", err)
		return gnet.Close
	}
	ctx, cancel := context.WithTimeout(wss.ctx, time.Second*5)
	defer cancel()

	//发送到msg服务进行处理
	ret, err := wss.msgRpcClient.SendMessage(ctx, &pb.ImMsgRequest{
		ChatType: msgCmd.ChatType,
		FromId:   msgCmd.FromId,
		Message:  msgCmd.Message,
		MsgId:    msgCmd.MsgId,
		MsgType:  msgCmd.MsgType,
		ToId:     msgCmd.ToId,
		Ts:       msgCmd.Ts,
	})
	if err != nil {
		logging.Errorf("routerClient.SendMessage err; %v", err)
		cmd := command.ImMsgCommandResp{
			Type: command.COMMAND_RESP_TYPE_MSG_ACK,
			Code: command.COMMAND_RESP_CODE_SERVER_ERR,
			Msg:  "server error",
			Data: map[string]interface{}{
				"msgId": msgCmd.MsgId,
			},
		}
		marshal, _ := json.Marshal(cmd)
		wsutil.WriteServerMessage(c, ws.OpText, marshal)
		return gnet.None
	}

	cmd := command.ImMsgCommandResp{
		Type: ret.Type,
		Code: ret.Code,
		Msg:  ret.Msg,
		Data: map[string]interface{}{
			"msgId":          ret.Data.MsgId,
			"sequence":       ret.Data.Sequence,
			"conversationId": ret.Data.ConversationId,
		},
	}
	marshal, _ := json.Marshal(cmd)
	wsutil.WriteServerMessage(c, ws.OpText, marshal)
	return gnet.None
}

func (wss *WsServer) OnTick() (delay time.Duration, action gnet.Action) {
	logging.Infof("[connected-count=%v]", atomic.LoadInt64(&wss.connected))
	return 10 * time.Second, gnet.None
}

// 推送消息给客户端
func (wss *WsServer) PushMsg(ctx context.Context, req *pb.PushRequest) (*emptypb.Empty, error) {
	logging.Infof("push to to:%d", req.UserId)
	var err error
	var once sync.Once
	wss.users.Range(func(key, value interface{}) bool {
		user := value.(*OnlineUser)
		if user.UserId == req.UserId {
			once.Do(func() {
				err = wsutil.WriteServerMessage(user.Conn, ws.OpText, []byte(req.Data))
			})
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// 客户端请求消息同步
func (ws *WsServer) OnMsgSync(c gnet.Conn, body []byte) gnet.Action {
	//todo 客户端同步消息
	//客户端发送最后一次的同步序列+会话ID
	return gnet.None
}
