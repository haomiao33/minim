package handler

import (
	"github.com/gofiber/fiber/v2"
	"im/internal/service/api/handler/msg"
	"im/internal/service/api/handler/recent_session"
)

var MsgHandler *msg.MsgHandler
var RecentSessionHandler *recent_session.RecentSessionHandler

func Init(router fiber.Router) {
	MsgHandler = msg.NewMsgHandler(router)
	RecentSessionHandler = recent_session.NewRecentSessionHandler(router)
}
