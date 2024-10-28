package msg

import (
	"errors"
	"github.com/guregu/null/v5"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gorm.io/gorm"
	"im/internal/model"
	"time"
)

type ImMsgDao struct {
}

func NewImMsgDao() *ImMsgDao {
	return &ImMsgDao{}
}

func (m *ImMsgDao) GetMsg(tx *gorm.DB, msgId string) (*model.ImMsg, error) {
	var msgModel model.ImMsg
	ret := tx.Table("im_msg").
		Where("id = ?", msgId).
		First(&msgModel)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &msgModel, nil
}

// 获取会话中大于指定序号的所有消息
func (m *ImMsgDao) GetMsgList(tx *gorm.DB, conversationId int64, sequence int64) ([]model.ImMsg, error) {
	var items []model.ImMsg
	ret := tx.Table("im_msg").
		Where("conversation_id = ?", conversationId).
		Where("sequence > ?", sequence).
		Where("status = ?", 0).
		Order("sequence asc").
		Find(&items)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return items, nil
}

func (m *ImMsgDao) AddMsg(tx *gorm.DB,
	conversationId int64,
	sequence int64,
	msgId string,
	chatType int,
	msgType int,
	fromId int64,
	toId int64,
	content string,
	msgTs int64,
	status int) (*model.ImMsg, error) {
	msg := &model.ImMsg{
		ConversationID: conversationId,
		Sequence:       sequence,
		ID:             msgId,
		ChatType:       chatType,
		MsgType:        msgType,
		FromID:         fromId,
		ToID:           toId,
		Content:        content,
		MsgTime:        null.TimeFrom(time.UnixMilli(msgTs)),
		Status:         status,
		CreatedTime:    null.TimeFrom(time.Now()),
		UpdateTime:     null.TimeFrom(time.Now()),
	}
	ret := tx.Table("im_msg").Create(msg)
	if ret.Error != nil {
		logging.Errorf("--- add msg failed:%v %d %d---",
			chatType, fromId, toId)
		return nil, ret.Error
	}
	return msg, nil
}
