package rpcclient

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"im/internal/logger"
	"im/pb"
)

func NewOffLineRpcClient(ctx context.Context, address string, serviceName string) pb.OfflineServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s/%s?wait=14s", address, serviceName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		logger.Fatalf("failed NewOffLineRpcClient err: %v", err)
	}
	return pb.NewOfflineServiceClient(conn)
}
