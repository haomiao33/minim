package main

import (
	"context"
	"im/internal/logger"
	"im/internal/service/msg_push"
	"im/internal/service/msg_push/config"
)

func main() {
	config.Init()

	logger.Init(config.Config.Log.Path,
		config.Config.Log.Level)
	defer logger.Sync()

	srv := msg_push.NewMsgPushServer(context.Background())
	srv.Run()
}
