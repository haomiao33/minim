package command

const (
	COMMAND_TYPE_HEARTBEAT = "heartbeat"
	COMMAND_TYPE_LOGIN_REQ = "login"
	COMMAND_TYPE_LOGIN_ACK = "loginAck"
	COMMAND_TYPE_LOGOUT    = "logout"

	//单聊消息
	COMMAND_TYPE_MSG = "msg"
	//单聊返回类型
	COMMAND_TYPE_MSG_ACK = "msgAck"

	//单聊消息-通知有新消息了sync notify
	COMMAND_TYPE_MSG_SYNC_NOTIFY = "msgSyncNotify"
	//单聊消息-同步单聊消息
	COMMAND_TYPE_MSG_SYNC = "msgSync"
	//单聊消息-同步单聊消息返回
	COMMAND_TYPE_MSG_SYNC_ACK = "msgSyncAck"
)

type ImCommand struct {
	//命令类型：
	//heartbeat 心跳
	//login 登录
	//logout 退出
	//im 消息
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
