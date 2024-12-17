package offline_push

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/modood/pushapi/huaweipush"
	"github.com/modood/pushapi/oppopush"
	"github.com/modood/pushapi/vivopush"
	"github.com/modood/pushapi/xiaomipush"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"im/internal/dao"
	"im/internal/db"
	"im/internal/logger"
	"im/internal/service/offline_push/config"
	"im/pb"
	"net"
	"strconv"
	"time"
)

const (
	VIVO_CONFIG   = "intent://com.xiaogongqiu.app/message?payload=%d#Intent;scheme=com.xiaogongqiu.app;launchFlags=0x10000000;component=com.xiaogongqiu.app/com.xiaogongqiu.app.MainActivity;end"
	OPPO_CONFIG   = "com.xiaogongqiu.app://message?payload=%d"
	HUAWEI_CONFIG = "intent://com.xiaogongqiu.app/message?payload=%d#Intent;scheme=com.xiaogongqiu.app;launchFlags=0x10000000;component=com.xiaogongqiu.app/com.xiaogongqiu.app.MainActivity;end"
	XIAOMI_CONFIG = "intent://com.xiaogongqiu.app/message?payload=%d#Intent;scheme=com.xiaogongqiu.app;launchFlags=0x10000000;component=com.xiaogongqiu.app/com.xiaogongqiu.app.MainActivity;end"
)

type OffLineServer struct {
	pb.UnimplementedOfflineServiceServer
	ctx context.Context

	oppoPushClient   *oppopush.Client
	vivoPushClient   *vivopush.Client
	huaweiPushClient *huaweipush.Client
	miPushClient     *xiaomipush.Client
}

func NewOffLineServer(ctx context.Context) *OffLineServer {
	oppoClient := oppopush.NewClient(config.Config.Oppo.AppKey,
		config.Config.Oppo.AppServerSecret)

	vivoClient := vivopush.NewClient(
		config.Config.Vivo.AppId,
		config.Config.Vivo.AppKey,
		config.Config.Vivo.AppSecret)

	huaweiClient := huaweipush.NewClient(config.Config.HuaWei.OAuthClientId,
		config.Config.HuaWei.OAuthClientSecret)

	miClient := xiaomipush.NewClient(config.Config.XiaoMi.AppSecret)

	return &OffLineServer{
		oppoPushClient:   oppoClient,
		vivoPushClient:   vivoClient,
		huaweiPushClient: huaweiClient,
		miPushClient:     miClient,
		ctx:              ctx,
	}
}

func (s *OffLineServer) oppoPush(toId int64, regId string, title string, content string) error {
	logger.Infof("oppo push data start toId: %d, regId:%s", toId, regId)

	var data = fmt.Sprintf(OPPO_CONFIG, toId)
	var req = &oppopush.SendReq{
		TargetType:  2,
		TargetValue: regId,
		Notification: &oppopush.Notification{
			Title:           title,
			Content:         content,
			ClickActionType: 5,
			ClickActionURL:  data,
		},
		VerifyRegistrationId: false,
	}
	send, err := s.oppoPushClient.Send(req)
	if err != nil {
		logger.Errorf("send oppo push failed:req:%+v, err: %+v", req, err)
		return err
	}
	logger.Infof("send oppo push result: %+v", send)
	return nil
}

func (s *OffLineServer) vivoPush(toId int64, regId string, title string, content string) error {
	logger.Infof("vivo push data start toId: %d, regId:%s", toId, regId)

	var data = fmt.Sprintf(VIVO_CONFIG, toId)
	var req = &vivopush.SendReq{
		RegId:          regId,
		NotifyType:     4,
		Title:          title,
		Content:        content,
		TimeToLive:     24 * 60 * 60,
		SkipType:       4,
		NetworkType:    -1,
		SkipContent:    data,
		Classification: 1,
		RequestId:      strconv.Itoa(int(time.Now().UnixNano())),
	}
	send, err := s.vivoPushClient.Send(req)
	if err != nil {
		logger.Errorf("send vivo push failed:req:%+v, err: %+v", req, err)
		return err
	}
	logger.Infof("send vivo push result: %+v", send)
	return nil
}

func (s *OffLineServer) huaweiPush(toId int64, regId string, title string, content string) error {
	logger.Infof("huawei push data start toId: %d, regId:%s", toId, regId)

	var data = fmt.Sprintf(HUAWEI_CONFIG, toId)
	var req = &huaweipush.SendReq{
		ValidateOnly: false,
		Message: &huaweipush.Message{
			Android: &huaweipush.AndroidConfig{
				Notification: &huaweipush.AndroidNotification{
					Title: title,
					Body:  content,
					ClickAction: &huaweipush.ClickAction{
						Type:   1,
						Intent: data,
					},
				},
			},
			Tokens: []string{regId},
		},
	}
	send, err := s.huaweiPushClient.Send(req)
	if err != nil {
		logger.Errorf("send huawei push failed:req:%+v, err: %+v", req, err)
		return err
	}
	logger.Infof("send huawei push result: %+v", send)
	return nil
}

func (s *OffLineServer) miPush(toId int64, regId string, title string, content string) error {
	logger.Infof("xiaomi push data start toId: %d, regId:%s", toId, regId)

	var data = fmt.Sprintf(XIAOMI_CONFIG, toId)
	var req = &xiaomipush.SendReq{
		RegistrationId: regId,
		Title:          title,
		Description:    content,
		NotifyType:     -1,
		Extra: &xiaomipush.Extra{
			NotifyEffect: "2",
			IntentUri:    data,
			ChannelId:    "130608",
		},
	}
	send, err := s.miPushClient.Send(req)
	if err != nil {
		logger.Errorf("send xiaomi push failed:req:%+v, err: %+v", req, err)
		return err
	}
	logger.Infof("send xiaomi push result: %+v", send)
	return nil
}

func (s *OffLineServer) Push(ctx context.Context, push *pb.OfflinePushRequest) (*emptypb.Empty, error) {
	user, err := dao.OffLineUserDao.GetOffLineUser(db.Db, push.UserId)
	if err != nil {
		logger.Warnf("------------get offline user  failed: %+v", err)
		return nil, err
	}
	if user.Platform == "oppo" {
		s.oppoPush(push.ConversationId, user.RegID, push.Title, push.Content)
	} else if user.Platform == "vivo" {
		s.vivoPush(push.ConversationId, user.RegID, push.Title, push.Content)
	} else if user.Platform == "huawei" {
		s.huaweiPush(push.ConversationId, user.RegID, push.Title, push.Content)
	} else if user.Platform == "mi" {
		s.miPush(push.ConversationId, user.RegID, push.Title, push.Content)
	} else {
		logger.Warnf("------------unknown platform: %s", user.Platform)
	}

	return &emptypb.Empty{}, nil
}

func (s *OffLineServer) Run() {
	cfg := config.Config.Rpc
	addr := fmt.Sprintf("%s:%d",
		cfg.ListenHost,
		cfg.ListenPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	logger.Infof("starting gRPC server on %s", addr)
	grpcServer := grpc.NewServer()
	pb.RegisterOfflineServiceServer(grpcServer, s)

	// 创建 Consul 客户端
	consulCfg := api.DefaultConfig()
	consulCfg.Address = config.Config.Consul.Address
	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		logger.Fatalf("failed to create consul client: %v", err)
	}

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// 创建服务注册信息
	//randomId, _ := uuid.NewUUID()
	registration := &api.AgentServiceRegistration{
		ID:      "OffLineService-1",
		Name:    "OffLineService",
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
		logger.Fatalf("failed to register service: %v", err)
	}

	logger.Infof("OffLine gRPC server is running on: %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}

	grpcServer.GracefulStop()
	logger.Info(" server exit")
}
