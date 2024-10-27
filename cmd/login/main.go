package main

import (
	"context"
	"im/internal/logger"
	"im/internal/service/login/config"
	"im/internal/service/login/ws"
)

func main() {
	config.Init()

	logger.Init(config.Config.Log.Path,
		config.Config.Log.Level)
	defer logger.Sync()

	loginServer := ws.NewWsServer(context.Background())
	loginServer.Run()
}
