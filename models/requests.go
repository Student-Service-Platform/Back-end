package models

import "gorm.io/gorm"

// Request
type Request struct {
	gorm.Model
	RequestID    int64   `json:"id"`                   // 反馈编号，ID 编号
	UserID       string  `json:"user_id"`              // 发起人id
	Title        string  `json:"title"`                // 反馈标题，反馈标题
	Discription  string  `json:"discription"`          // 反馈内容
	Category     *int64  `json:"category,omitempty"`   // 反馈问题类别，标签
	Urgency      int64   `json:"urgency"`              // 紧急程度，1-5，数值越大紧急程度越高
	UndertakerID *string `json:"undertaker,omitempty"` // 承接人，接单的管理员id
	IfRubbish    int64   `json:"if_rubbish"`           // 状态，默认1,计入评价状况
	IsAnonymous  bool    `json:"is_anonymous"`         // 是否匿名，默认false（不匿名）
	Status       bool    `json:"status"`               // 是否已经处理
	Grade        int     `json:"grade"`                // 评分
	GradeContent string  `json:"grade_content"`        // 评价内容
}
