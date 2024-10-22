package rpcclient

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"google.golang.org/grpc"
	"im/pb"
)

func NewMsgRpcClient(ctx context.Context, address string, serviceName string) pb.MsgServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s/%s?wait=14s", address, serviceName),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		logging.Fatalf("failed NewMsgServiceClient err: %v", err)
	}
	return pb.NewMsgServiceClient(conn)
}
