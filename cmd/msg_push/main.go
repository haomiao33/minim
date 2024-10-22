package main

import (
	"context"
	"im/internal/service/msg_push"
)

func main() {
	srv := msg_push.NewMsgPushServer(context.Background())
	srv.Run()
}
