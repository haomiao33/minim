package handler

import (
	"github.com/gofiber/fiber/v2"
	"im/internal/service/api/handler/msg"
)

var MsgHandler msg.MsgHandler

func Init(router fiber.Router) {
	MsgHandler = msg.NewMsgHandler(router)
}
