package msg

import (
	"errors"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gorm.io/gorm"
	"im/internal/model"
	"im/internal/sharding"
	"time"
)

type ImMsgDao struct {
}

func NewImMsgDao() *ImMsgDao {
	return &ImMsgDao{}
}

func (m *ImMsgDao) GetMsg(tx *gorm.DB, msgId string, conversationId int64) (*model.ImMsg, error) {
	var msgModel model.ImMsg
	ret := tx.Table(sharding.GetTableName("im_msg", conversationId)).
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
	ret := tx.Table(sharding.GetTableName("im_msg", conversationId)).
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
		MsgTime:        model.MyTime{time.UnixMilli(msgTs)},
		Status:         status,
		CreatedTime:    model.MyTime{time.Now()},
		UpdatedTime:    model.MyTime{time.Now()},
	}
	ret := tx.Table(sharding.GetTableName("im_msg", conversationId)).Create(msg)
	if ret.Error != nil {
		logging.Errorf("--- add msg failed:%v %d %d---",
			chatType, fromId, toId)
		return nil, ret.Error
	}
	return msg, nil
}

func (m *ImMsgDao) UpdateMsg(tx *gorm.DB, conversationId int64, msgIds []string, body map[string]interface{}) error {
	ret := tx.Table(sharding.GetTableName("im_msg", conversationId)).
		Where("id in ?", msgIds).
		Updates(body)
	if ret.Error != nil {
		logging.Errorf("--- UpdateMsg failed:%v, ---", msgIds)
		return ret.Error
	}
	return nil
}
