package model

import "github.com/guregu/null/v5"

// ImMsg [...]
type ImMsg struct {
	ID             string      `gorm:"primaryKey;column:id" json:"id"`
	ConversationID int64       `gorm:"column:conversation_id" json:"conversationId"`
	MsgType        int         `gorm:"column:msg_type" json:"msgType"` // 消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
	FromID         int64       `gorm:"column:from_id" json:"fromId"`
	ToID           int64       `gorm:"column:to_id" json:"toId"`
	ChatType       int         `gorm:"column:chat_type" json:"chatType"` // 0=单聊；1=一般群； 2=机器人
	Content        string      `gorm:"column:content" json:"content"`
	Status         int         `gorm:"column:status" json:"status"` // 0=已发送, 1=已送达, 2=已读, 3=已撤回
	MsgRead        bool        `gorm:"column:msg_read" json:"msgRead"`
	Sequence       int64       `gorm:"column:sequence" json:"sequence"` // 消息顺序
	ReplyTo        int64       `gorm:"column:reply_to" json:"replyTo"`
	MsgAudit       int         `gorm:"column:msg_audit" json:"msgAudit"` // 0=默认
	RefID          null.String `gorm:"column:ref_id" json:"refId"`       // 关联消息id
	Revoked        bool        `gorm:"column:revoked" json:"revoked"`
	MsgTime        null.Time   `gorm:"column:msg_time" json:"msgTime"`
	RevokedTime    null.Time   `gorm:"column:revoked_time" json:"revokedTime"`
	RevokedBy      int64       `gorm:"column:revoked_by" json:"revokedBy"`
	CreatedTime    null.Time   `gorm:"column:created_time" json:"createdTime"`
	UpdateTime     null.Time   `gorm:"column:update_time" json:"updateTime"`
	DeleteTime     null.Time   `gorm:"column:delete_time" json:"deleteTime"`
}

// TableName get sql table name.获取数据库表名
func (m *ImMsg) TableName() string {
	return "im_msg"
}
