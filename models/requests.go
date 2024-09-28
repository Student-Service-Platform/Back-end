package models

import "gorm.io/gorm"

// Request
type Request struct {
	gorm.Model
	// 反馈问题类别，标签，1土木类，2水电气类，3公共设施问题，4其它问题
	Category *int64 `json:"category,omitempty"`
	// 反馈内容
	Discription string `json:"discription"`
	// 反馈编号，ID 编号
	ID int64 `json:"id"`
	// 是否匿名，默认false（不匿名）
	IsAnonymous bool `json:"is_anonymous"`
	// 状态，默认1,计入评价状况
	Status int64 `json:"status"`
	// 反馈标题，反馈标题
	Title string `json:"title"`
	// 承接人，承接管理员id
	Undertaker *string `json:"undertaker,omitempty"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
	// 紧急程度，1-5，数值越大紧急程度越高
	Urgency int64 `json:"urgency"`
	// 发起人id
	UserID string `json:"user_id"`
}
