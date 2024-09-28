package models

import "gorm.io/gorm"

// Reply 回复
type Reply struct {
	gorm.Model
	// 回复内容
	Content string `json:"content"`
	// 创建时间（gorm管）
	// 回复id，ID 编号
	ID int64 `json:"id"`
	// 反馈id，回复的反馈的id
	RequestID *int64 `json:"request_id,omitempty"`
	// 回复者id
	Respondent string `json:"respondent"`
	// 回复的是回复还是反馈
	Type string `json:"type"`
}
