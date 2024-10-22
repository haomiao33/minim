package command

// 登陆消息
type ImLoginCommandReq struct {
	UserId int64 `json:"userId"`
}

// 登陆返回
type ImLoginCommandResp struct {
	Type string `json:"type"`
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
