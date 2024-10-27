package resp

type MsgSendResp struct {
	MsgId          string `json:"msgId"`
	Sequence       int64  `json:"sequence"`
	ConversationId int64  `json:"conversationId"`
}
