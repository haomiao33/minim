package main

import (
	"context"
	"im/internal/db"
	"im/internal/logger"
	"im/internal/service/offline_push"
	"im/internal/service/offline_push/config"
)

func main() {
	config.Init()

	logger.Init(config.Config.Log.Path,
		config.Config.Log.Level)
	defer logger.Sync()

	cfg := config.Config
	db.Init(cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database)

	srv := offline_push.NewOffLineServer(context.Background())
	srv.Run()
}
