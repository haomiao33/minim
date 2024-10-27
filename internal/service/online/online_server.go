package online

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"im/internal/common"
	"im/internal/logger"
	"im/internal/service/online/config"
	"im/pb"
	"im/pkg/redisclient"
	"net"
	"time"
)

type OnlineServer struct {
	pb.UnimplementedOnlineServiceServer
	redisClient *redisclient.RedisClient
	ctx         context.Context
}

func NewOnlineServer(ctx context.Context) *OnlineServer {
	return &OnlineServer{
		ctx: ctx,
		redisClient: redisclient.NewRedisClient(ctx,
			config.Config.Redis.Addr,
			config.Config.Redis.Password),
	}
}

func (s *OnlineServer) UpdateOnlineUser(ctx context.Context,
	req *pb.UpdateOnlineUserRequest) (*emptypb.Empty, error) {
	user := &common.OnlineUser{
		UserId:       req.UserId,
		ServerId:     req.ServerId,
		LastUpdateTs: req.LastUpdateTs,
	}
	marshal, _ := json.Marshal(user)
	s.redisClient.Client.Set(ctx,
		fmt.Sprintf("im:online:%d", req.UserId),
		string(marshal), time.Duration(config.Config.App.OnlineUserTimeOutSeconds)*time.Second)
	return &emptypb.Empty{}, nil
}

func (s *OnlineServer) GetOnlineUser(ctx context.Context,
	req *pb.GetOnlineUserRequest) (*pb.GetOnlineUserResp, error) {
	get := s.redisClient.Client.Get(ctx, fmt.Sprintf("im:online:%d", req.UserId))
	if get.Err() != nil {
		return nil, get.Err()
	}
	if get.Val() == "" {
		return nil, fmt.Errorf("user not online")
	}
	var user common.OnlineUser
	err := json.Unmarshal([]byte(get.Val()), &user)
	if err != nil {
		return nil, get.Err()
	}

	return &pb.GetOnlineUserResp{
		UserId:       user.UserId,
		ServerId:     user.ServerId,
		LastUpdateTs: user.LastUpdateTs,
	}, nil
}

func (s *OnlineServer) OutlineUser(ctx context.Context,
	req *pb.GetOnlineUserRequest) (*emptypb.Empty, error) {
	s.redisClient.Client.Del(ctx, fmt.Sprintf("im:online:%d", req.UserId))
	return &emptypb.Empty{}, nil
}

func (s *OnlineServer) Run() {
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
	pb.RegisterOnlineServiceServer(grpcServer, s)

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
		ID:      "OnlineService-1",
		Name:    "OnlineService",
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

	logger.Infof("online gRPC server is running on: %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}

	grpcServer.GracefulStop()
	logger.Info(" server exit")
}
