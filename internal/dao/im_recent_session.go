package dao

import (
	"errors"
	"github.com/guregu/null/v5"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gorm.io/gorm"
	"im/internal/model"
	"im/internal/service/api/resp"
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

// 获取某个用户的会话列表
func (s *ImRecentSessionDao) GetConversationList(tx *gorm.DB,
	userId int64, chatType int32) ([]resp.ImUserRecentSessionResp, error) {
	items := make([]resp.ImUserRecentSessionResp, 0)
	ret := tx.Table("im_recent_session a").
		Joins("left join im_conversation b on a.conversation_id = b.id").
		Select("a.*,b.sequence").
		Where("a.user_id = ?", userId).
		Where("a.type = ?", chatType).
		Where("a.deleted_time IS NULL").
		Order("a.updated_time DESC,b.sequence desc").
		Find(&items)
	if ret.Error != nil {
		logging.Errorf("--- GetConversationList conversation failed:%d, ---", userId)
		return items, nil
	}
	return items, nil
}

// 删除会话
func (s *ImRecentSessionDao) DelConversation(tx *gorm.DB,
	conversationId int64, userId int64) error {
	ret := tx.Table("im_recent_session").
		Where("conversation_id = ? and user_id = ?", conversationId, userId).
		Updates(map[string]interface{}{
			"deleted_time": time.Now(),
		})
	if ret.Error != nil {
		logging.Errorf("--- DelConversation conversation failed:%d, ---", userId)
		return ret.Error
	}
	return nil
}
