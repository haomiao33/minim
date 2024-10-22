package main

import (
	"context"
	"im/internal/service/online"
)

func main() {
	srv := online.NewOnlineServer(context.Background())
	srv.Run()
}
