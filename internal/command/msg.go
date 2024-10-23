package command

const (
	COMMAND_RESP_CODE_SUCESS     = 10000
	COMMAND_RESP_CODE_SERVER_ERR = 10001
)

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

// 单聊ack
type ImMsgCommandResp struct {
	//ack
	Type string      `json:"type"`
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 单聊sync通知
type ImMsgSyncNotifyCommand struct {
	ChatType int
	MsgType  int
	MsgId    string
	//会话id
	ConversationId int64
	//消息序号
	Sequence int64
	FromId   int64
	ToId     int64
}

// 同步单聊消息
type ImMsgSyncCommandReq struct {
	ConversationId int64 `json:"conversationId"`
	Sequence       int64 `json:"sequence"`
	UserId         int64 `json:"userId"`
}

// 同步单聊消息返回
// 单聊 消息 item
type ImMsgItem struct {
	Id             string `json:"id"`
	ConversationId int64  `json:"conversationId"`
	MsgType        int    `json:"msgType"`
	FromId         int64  `json:"fromId"`
	ToId           int64  `json:"toId"`
	ChatType       int    `json:"chatType"`
	Content        string `json:"content"`
	Status         int    `json:"status"`
	MsgRead        int    `json:"msgRead"`
	Sequence       int64  `json:"sequence"`
	ReplyTo        int64  `json:"replyTo"`
	MsgAudit       int    `json:"msgAudit"`
	RefId          string `json:"refId"`
	Revoked        int    `json:"revoked"`
	RevokedTime    int64  `json:"revokedTime"`
	RevokedBy      int64  `json:"revokedBy"`
	CreatedTime    int64  `json:"createdTime"`
	UpdateTime     int64  `json:"updateTime"`
}
