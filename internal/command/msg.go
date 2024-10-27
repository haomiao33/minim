package command

const (
	//单聊消息-通知有新消息了sync notify
	COMMAND_TYPE_MSG_SYNC_NOTIFY = "msgSyncNotify"
)

// 单聊sync通知
type ImMsgSyncNotifyCommand struct {
	ChatType int    `json:"chatType"`
	MsgType  int    `json:"msgType"`
	MsgId    string `json:"msgId"`
	//会话id
	ConversationId int64 `json:"conversationId"`
	//消息序号
	Sequence int64 `json:"sequence"`
	FromId   int64 `json:"fromId"`
	ToId     int64 `json:"toId"`
}
