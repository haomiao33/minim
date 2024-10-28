package user

import (
	"errors"
	"gorm.io/gorm"
	"im/internal/model"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (m *UserDao) GetUser(tx *gorm.DB, userId int64) (*model.User, error) {
	var userModel model.User
	ret := tx.Table("user").
		Where("user_id = ?", userId).
		First(&userModel)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &userModel, nil
}

func (m *UserDao) GetUserByFiled(tx *gorm.DB, userId int64, fields []string) (*model.User, error) {
	var userModel model.User
	ret := tx.Table("user").
		Where("user_id = ?", userId).
		Select(fields).
		First(&userModel)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &userModel, nil
}
