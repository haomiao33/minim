package model

// OfflinePushUser [...]
type OfflinePushUser struct {
	ID         int64  `gorm:"primaryKey;column:id" json:"id"`
	UserID     int64  `gorm:"column:user_id" json:"userId"`
	Platform   string `gorm:"column:platform" json:"platform"`      // 用户所属平台：oppo、vivo、mi、huawei、apple
	RegID      string `gorm:"column:reg_id" json:"regId"`           // 设备注册id
	CreateTime MyTime `gorm:"column:create_time" json:"createTime"` // 创建时间
	UpdateBy   int64  `gorm:"column:update_by" json:"updateBy"`     // 修改者
	UpdateTime MyTime `gorm:"column:update_time" json:"updateTime"` // 修改时间
	CreateBy   int64  `gorm:"column:create_by" json:"createBy"`     // 创建者
}

// TableName get sql table name.获取数据库表名
func (m *OfflinePushUser) TableName() string {
	return "offline_push_user"
}
