package main

import (
	"context"
	"im/internal/service/login/ws"
)

func main() {
	loginServer := ws.NewWsServer(context.Background())
	loginServer.Run()
}
