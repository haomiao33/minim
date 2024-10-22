package dao

var MsgDao *ImMsgDao
var ConversationDao *ImConversationDao
var RecentSessionDao *ImRecentSessionDao

func init() {
	MsgDao = NewImMsgDao()
	ConversationDao = NewConversationDao()
	RecentSessionDao = NewRecentSessionDao()
}
