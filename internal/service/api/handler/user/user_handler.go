package user

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"im/internal/db"
	"im/internal/model"
	"im/internal/response"
	"im/internal/service/api/req"
)

type UserHandler struct {
}

func NewUserHandler(router fiber.Router) *UserHandler {
	handler := &UserHandler{}
	// 设备注册
	router.Post("/user/device/register", handler.DeviceRegister)
	return handler
}

func (u *UserHandler) DeviceRegister(c *fiber.Ctx) error {
	var req req.DeviceRegisterReq
	if err := c.BodyParser(&req); err != nil {
		return errors.New("参数错误")
	}
	var modelUser model.OfflinePushUser
	tx := db.Db.Model(&model.OfflinePushUser{}).First(&modelUser,
		"user_id = ? and platform = ?",
		req.UserId, req.Platform)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return errors.New("出错啦")
		}
		//create
		modelUser = model.OfflinePushUser{
			UserID:   req.UserId,
			Platform: req.Platform,
			RegID:    req.RegId,
		}
		db.Db.Save(&modelUser)
	} else {
		if modelUser.RegID != req.RegId {
			//update
			db.Db.Model(&model.OfflinePushUser{}).
				Where("user_id = ? and platform = ?",
					req.UserId, req.Platform).
				Update("reg_id", req.RegId)
		}
	}

	return c.JSON(response.Success(""))
}
