package rpcclient

import (
	"context"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"google.golang.org/grpc"
	"im/pb"
)

func NewLoginDefaultRpcClient(ctx context.Context, address string, serviceName string) pb.LoginServiceClient {
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
	)
	if err != nil {
		logging.Fatalf("failed NewLoginDefaultRpcClient err: %v", err)
	}
	return pb.NewLoginServiceClient(conn)
}
