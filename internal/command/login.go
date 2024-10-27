package command

const (
	COMMAND_TYPE_HEARTBEAT = "heartbeat"
	COMMAND_TYPE_LOGIN_REQ = "login"
	COMMAND_TYPE_LOGIN_ACK = "loginAck"
	COMMAND_TYPE_LOGOUT    = "logout"
)

// 登陆消息
type ImLoginCommandReq struct {
	UserId int64 `json:"userId"`
}
