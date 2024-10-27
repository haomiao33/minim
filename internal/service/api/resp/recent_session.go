package resp

import "im/internal/model"

// 获取用户会话列表resp
type ImUserRecentSessionResp struct {
	model.ImRecentSession
	Sequence int64 `json:"sequence"`
}
