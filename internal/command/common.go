package command

// websocket command req
type ImCommandReq struct {
	//命令类型：
	//heartbeat 心跳
	//login 登录
	//logout 退出
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// websocket command resp
type ImCommandResp struct {
	//命令类型：
	//heartbeat 心跳
	//login 登录
	//logout 退出
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Code    int32       `json:"code"`
	Message string      `json:"message"`
}
