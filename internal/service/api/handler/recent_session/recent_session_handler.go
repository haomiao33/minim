package recent_session

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"im/internal/dao"
	"im/internal/db"
	"im/internal/response"
	"im/internal/service/api/req"
)

type RecentSessionHandler struct {
}

func NewRecentSessionHandler(router fiber.Router) *RecentSessionHandler {
	handler := &RecentSessionHandler{}
	router.Post("/RecentSession/userRecentSessionList", handler.GetUserRecentSession)
	router.Post("/RecentSession/deleteUserSession", handler.DeleteUserSession)
	return handler
}

// GetUserRecentSession 获取用户最近会话列表
func (u *RecentSessionHandler) GetUserRecentSession(c *fiber.Ctx) error {
	var msg req.ImUserRecentSessionReq
	if err := c.BodyParser(&msg); err != nil {
		return errors.New("参数错误")
	}
	list, _ := dao.RecentSessionDao.GetConversationList(db.Db, msg.UserId, msg.ChatType)

	return c.JSON(response.Success(list))
}

// GetUserRecentSession 获取用户最近会话列表
func (u *RecentSessionHandler) DeleteUserSession(c *fiber.Ctx) error {
	var msg req.ImDeleteUserRecentSessionReq
	if err := c.BodyParser(&msg); err != nil {
		return errors.New("参数错误")
	}
	err := dao.RecentSessionDao.DelConversation(db.Db, msg.ConversationId, msg.UserId)
	if err != nil {
		return errors.New("删除失败")
	}
	return c.JSON(response.Success("success"))
}
