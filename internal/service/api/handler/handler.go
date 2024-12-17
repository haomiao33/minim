package handler

import (
	"github.com/gofiber/fiber/v2"
	"im/internal/service/api/handler/msg"
	"im/internal/service/api/handler/recent_session"
	"im/internal/service/api/handler/upload"
	"im/internal/service/api/handler/user"
)

var MsgHandler *msg.MsgHandler
var RecentSessionHandler *recent_session.RecentSessionHandler
var UploadHandler *upload.UploadHandler
var UserHandler *user.UserHandler

func Init(router fiber.Router) {
	MsgHandler = msg.NewMsgHandler(router)
	RecentSessionHandler = recent_session.NewRecentSessionHandler(router)
	UploadHandler = upload.NewUploadHandler(router)
	UserHandler = user.NewUserHandler(router)
}
