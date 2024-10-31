package model

// ImConversation [...]
type ImConversation struct {
	ID          int64  `gorm:"primaryKey;column:id" json:"id"`
	SmallID     int64  `gorm:"column:small_id" json:"smallId"`
	BigID       int64  `gorm:"column:big_id" json:"bigId"`
	Type        int    `gorm:"column:type" json:"type"`         // 0=单聊；1=一般群； 2=机器人
	Sequence    int64  `gorm:"column:sequence" json:"sequence"` // 消息顺序
	CreatedTime MyTime `gorm:"column:created_time" json:"createdTime"`
	UpdatedTime MyTime `gorm:"column:updated_time" json:"updatedTime"`
}

// TableName get sql table name.获取数据库表名
func (m *ImConversation) TableName() string {
	return "im_conversation"
}
