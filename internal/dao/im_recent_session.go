package dao

import (
	"errors"
	"github.com/guregu/null/v5"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gorm.io/gorm"
	"im/internal/model"
	"time"
)

type ImRecentSessionDao struct {
}

func NewRecentSessionDao() *ImRecentSessionDao {
	return &ImRecentSessionDao{}
}

func (s *ImRecentSessionDao) Get(tx *gorm.DB, chatType int32, userId int64, otherId int64) (*model.ImRecentSession, error) {
	var session model.ImRecentSession
	ret := tx.Table("im_recent_session").
		Where("type = ? and user_id = ? and other_id = ?", chatType, userId, otherId).
		First(&session)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &session, nil
}

func (s *ImRecentSessionDao) Add(tx *gorm.DB, conversationId int64, chatType int32, userId int64, otherId int64,
	lastMsgId string, lastMsg string, lastMsgTime time.Time) error {
	session := model.ImRecentSession{
		UserID:         userId,
		OtherID:        otherId,
		Type:           int(chatType),
		ConversationId: conversationId,
		LastMsgId:      lastMsgId,
		LastMsg:        lastMsg,
		LastMsgTime:    lastMsgTime,
		CreatedTime:    null.TimeFrom(time.Now()),
		UpdatedTime:    null.TimeFrom(time.Now()),
		SessionMute:    0,
		SessionTop:     0,
	}
	ret := tx.Table("im_recent_session").Create(&session)
	if ret.Error != nil {
		logging.Errorf("--- create recent session failed:%v %d %d---",
			chatType, userId, otherId)
		return ret.Error
	}
	return nil
}

func (s *ImRecentSessionDao) Update(tx *gorm.DB, chatType int32, userId int64, otherId int64,
	lastMsgId string, lastMsg string, lastMsgTime time.Time) error {
	ret := tx.Table("im_recent_session").
		Where("type = ? and user_id = ? and other_id = ?", chatType, userId, otherId).
		Updates(map[string]interface{}{
			"last_msg_id":   lastMsgId,
			"last_msg":      lastMsg,
			"last_msg_time": lastMsgTime,
		})
	if ret.Error != nil {
		logging.Errorf("--- update recent session failed:%v %d %d---",
			chatType, userId, otherId)
		return ret.Error
	}
	return nil
}
