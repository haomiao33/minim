package conversation

import (
	"errors"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"gorm.io/gorm"
	"im/internal/model"
	"time"
)

type ImConversationDao struct {
}

func NewConversationDao() *ImConversationDao {
	return &ImConversationDao{}
}

// /获取会话
func (d *ImConversationDao) GetConversation(tx *gorm.DB, chatType int32, bigId int64, smallId int64) (*model.ImConversation, error) {
	var conversation model.ImConversation
	ret := tx.Table("im_conversation").
		Where("type = ? and big_id = ? and small_id = ?", chatType, bigId, smallId).
		First(&conversation)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return &conversation, nil
}

// 添加会话
func (d *ImConversationDao) AddConversation(tx *gorm.DB, chatType int32, bigId int64, smallId int64) (*model.ImConversation, error) {
	//创建会话关系
	conversation := model.ImConversation{
		BigID:       bigId,
		SmallID:     smallId,
		Type:        int(chatType),
		Sequence:    0,
		CreatedTime: model.MyTime{time.Now()},
		UpdatedTime: model.MyTime{time.Now()},
	}
	ret := tx.Table("im_conversation").Create(&conversation)
	if ret.Error != nil {
		logging.Errorf("--- create conversation failed:%v %d %d---",
			chatType, bigId, smallId)
		return nil, ret.Error
	}
	return &conversation, nil
}

// 更新会话
func (d *ImConversationDao) UpdateConversation(tx *gorm.DB, conversationId int64, sequence int64) error {
	ret := tx.Table("im_conversation").
		Where("id = ?", conversationId).
		Updates(map[string]interface{}{"sequence": sequence})
	if ret.Error != nil {
		logging.Errorf("--- update conversation failed:%d, seq:%d---", conversationId, sequence)
		return ret.Error
	}
	return nil
}
