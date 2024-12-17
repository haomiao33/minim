package offline_push_user

import (
	"gorm.io/gorm"
	"im/internal/model"
)

type OffLinePushUserDao struct {
}

func NewOffLinePushUserDao() *OffLinePushUserDao {
	return &OffLinePushUserDao{}
}

func (m *OffLinePushUserDao) GetOffLineUser(tx *gorm.DB, userId int64) (*model.OfflinePushUser, error) {
	var userModel model.OfflinePushUser
	ret := tx.Model(&model.OfflinePushUser{}).
		Where("user_id = ?", userId).
		First(&userModel)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &userModel, nil
}
