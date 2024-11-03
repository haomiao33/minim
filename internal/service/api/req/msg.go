package req

// 单聊消息req
type ImMsgCommandReq struct {
	//msgId
	MsgId string `json:"msgId"`
	//0=单聊；1=一般群； 2=机器人
	ChatType int32 `json:"chatType"`
	//消息类型； 1=文本；2=图片；3=视频；4=文件；5=通话
	MsgType int32  `json:"msgType"`
	FromId  int64  `json:"fromId"`
	ToId    int64  `json:"toId"`
	Content string `json:"content"`
	//消息时间戳
	Ts int64 `json:"ts"`
}

// 同步单聊消息
type ImMsgSyncCommandReq struct {
	//会话id
	ConversationId int64 `json:"conversationId"`
	//客户端存放的最新消息序号
	Sequence int64 `json:"sequence"`
	//用户id
	UserId int64 `json:"userId"`
	//otherId 如果为0 就不传递用户信息，大于0 就获取这个otherId的用户信息
	OtherId int64 `json:"otherId"`
}
