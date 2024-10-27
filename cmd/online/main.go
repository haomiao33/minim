package main

import (
	"context"
	"im/internal/logger"
	"im/internal/service/online"
	"im/internal/service/online/config"
)

func main() {
	config.Init()

	logger.Init(config.Config.Log.Path,
		config.Config.Log.Level)
	defer logger.Sync()

	srv := online.NewOnlineServer(context.Background())
	srv.Run()
}
