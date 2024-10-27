package rpcclient

import (
	"context"
	"google.golang.org/grpc"
	"im/internal/logger"
	"im/pb"
)

func NewLoginDefaultRpcClient(ctx context.Context, address string, serviceName string) pb.LoginServiceClient {
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Fatalf("failed NewLoginDefaultRpcClient err: %v", err)
	}
	return pb.NewLoginServiceClient(conn)
}
