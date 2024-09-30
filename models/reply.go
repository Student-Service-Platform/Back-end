package models

import "gorm.io/gorm"

// Reply 回复
type Reply struct {
	gorm.Model        // 创建时间（gorm管）
	ReplyID    int64  `json:"reply_id"`             // 回复id，ID 编号
	RequestID  *int64 `json:"request_id,omitempty"` // 反馈id，回复的反馈的id
	Father     int64  `json:"father"`               //小林给个注释看
	Content    string `json:"content"`              // 回复内容
	Respondent string `json:"respondent"`           // 回复者id

}
