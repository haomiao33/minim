package model

// ImRecentSession [...]
type ImRecentSession struct {
	UserID         int64  `gorm:"primaryKey;column:user_id" json:"userId"`
	OtherID        int64  `gorm:"primaryKey;column:other_id" json:"otherId"`
	ConversationId int64  `gorm:"column:conversation_id" json:"conversationId"` //会话id
	Type           int    `gorm:"primaryKey;column:type" json:"type"`           // 0=单聊；1=一般群； 2=机器人
	LastMsgId      string `gorm:"column:last_msg_id" json:"lastMsgId"`
	LastMsg        string `gorm:"column:last_msg" json:"lastMsg"`
	LastMsgTime    MyTime `gorm:"column:last_msg_time" json:"lastMsgTime"`
	CreatedTime    MyTime `gorm:"column:created_time" json:"createdTime"`
	UpdatedTime    MyTime `gorm:"column:updated_time" json:"updatedTime"`
	DeletedTime    MyTime `gorm:"column:deleted_time" json:"deletedTime"`
	SessionMute    int8   `gorm:"column:session_mute" json:"sessionMute"` // 是否禁止提醒；0=否；1=是
	SessionTop     int8   `gorm:"column:session_top" json:"sessionTop"`   // 是否置顶
}

// TableName get sql table name.获取数据库表名
func (m *ImRecentSession) TableName() string {
	return "im_recent_session"
}
