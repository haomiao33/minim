package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	_ "google.golang.org/grpc/health"
	"im/internal/service/api/config"
	"im/internal/service/api/handler"
	"im/internal/service/api/middleware"
)

type MsgPushServer struct {
	app *fiber.App
	ctx context.Context
}

func NewApiServer(ctx context.Context) *MsgPushServer {
	//fiber app
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})
	app.Use(recover2.New(recover2.Config{
		EnableStackTrace: true,
	}))
	//返回统一处理
	app.Use(middleware.ErrorMiddlewareResponse)
	group := app.Group("/api/v1/")
	handler.Init(group)

	return &MsgPushServer{
		app: app,
		ctx: ctx,
	}
}

func (s *MsgPushServer) Run() {
	log.Infof("server start at %s", config.Config.App.Listener)
	log.Fatal(s.app.Listen(config.Config.App.Listener))
	log.Info("Server exit")
}