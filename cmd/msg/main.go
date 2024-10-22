package main

import (
	"context"
	"im/internal/db"
	"im/internal/service/msg"
	"im/internal/service/msg/config"
)

func main() {
	cfg := config.Config
	db.Init(cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database)
	srv := msg.NewMsgServer(context.Background())
	srv.Run()
}
