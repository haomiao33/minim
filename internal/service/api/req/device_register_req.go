package req

// 设备注册请求
type DeviceRegisterReq struct {
	UserId   int64  `json:"userId"`
	RegId    string `json:"regId"`
	Platform string `json:"platform"`
}
