package common

type OnlineUser struct {
	UserId       int64  `json:"userId"`
	ServerId     string `json:"serverId"`
	LastUpdateTs int64  `json:"lastUpdateTs"`
}
