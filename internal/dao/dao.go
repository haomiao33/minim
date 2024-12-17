package dao

import (
	"im/internal/dao/conversation"
	"im/internal/dao/msg"
	"im/internal/dao/offline_push_user"
	"im/internal/dao/recent_session"
	"im/internal/dao/user"
)

var MsgDao *msg.ImMsgDao
var ConversationDao *conversation.ImConversationDao
var RecentSessionDao *recent_session.ImRecentSessionDao
var UserDao *user.UserDao
var OffLineUserDao *offline_push_user.OffLinePushUserDao

func init() {
	MsgDao = msg.NewImMsgDao()
	ConversationDao = conversation.NewConversationDao()
	RecentSessionDao = recent_session.NewRecentSessionDao()
	UserDao = user.NewUserDao()
	OffLineUserDao = offline_push_user.NewOffLinePushUserDao()
}
