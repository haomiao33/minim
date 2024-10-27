package req

// 单聊消息req
type ImMsgCommandReq struct {
	//msgId,chatType,msgType,from,to,message,ts
	MsgId string `json:"msgId"`
	//0=单聊；1=一般群； 2=机器人
	ChatType int32 `json:"chatType"`
	//消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
	MsgType int32  `json:"msgType"`
	FromId  int64  `json:"fromId"`
	ToId    int64  `json:"toId"`
	Message string `json:"message"`
	//消息时间戳
	Ts int64 `json:"ts"`
}

// 同步单聊消息
type ImMsgSyncCommandReq struct {
	ConversationId int64 `json:"conversationId"`
	Sequence       int64 `json:"sequence"`
	UserId         int64 `json:"userId"`
}
