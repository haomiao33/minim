package req

// 获取用户会话列表req
type ImUserRecentSessionReq struct {
	UserId int64 `json:"userId"`
	//0=单聊；1=一般群； 2=机器人
	ChatType int32 `json:"chatType"`
}

// 获取用户会话列表req
type ImDeleteUserRecentSessionReq struct {
	UserId int64 `json:"userId"`
	//会话id
	ConversationId int64 `json:"conversationId"`
}
