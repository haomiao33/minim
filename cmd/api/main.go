package main

import (
	"context"
	"im/internal/db"
	"im/internal/service/api"
	"im/internal/service/api/client"
	"im/internal/service/api/config"
)

func main() {
	ctx := context.Background()

	cfg := config.Config
	db.Init(cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database)
	client.Init(ctx)
	srv := api.NewApiServer(ctx)
	srv.Run()
}
