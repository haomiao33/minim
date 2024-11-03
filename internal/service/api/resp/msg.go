package resp

import (
	"im/internal/model"
)

type MsgSendResp struct {
	Id             string       `json:"id"`
	Sequence       int64        `json:"sequence"`
	Status         int          `json:"status"`
	ConversationId int64        `json:"conversationId"`
	RevokedTime    model.MyTime `json:"revokedTime"`
	CreatedTime    model.MyTime `json:"createdTime"`
	UpdatedTime    model.MyTime `json:"updatedTime"`
	DeletedTime    model.MyTime `json:"deletedTime"`
}

type ImUserInfoResp struct {
	UserID   int64  `gorm:"primaryKey;column:user_id" json:"userId"` // 用户Id
	NickName string `gorm:"column:nick_name" json:"nickName"`        // 用户昵称
	Avatar   string `gorm:"column:avatar" json:"avatar"`             // 头像地址
	UserType int    `gorm:"column:user_type" json:"userType"`        // 用户类型:user=0;admin=10000;机构=300
}

// 同步单聊消息返回
type ImMsgSyncCommandResp struct {
	Items     []model.ImMsg   `json:"items"`
	OtherInfo *ImUserInfoResp `json:"otherInfo"`
}
