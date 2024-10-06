package models

import "gorm.io/gorm"

// Request 反馈
type Request struct {
	gorm.Model
	UserID       string  `gorm:"type:varchar(100);not null"`
	Title        string  `json:"title"`              // 反馈标题，反馈标题
	Description  string  `json:"description"`        // 反馈内容
	Category     int     `json:"category,omitempty"` // 反馈问题类别，标签
	Urgency      int     `json:"urgency"`            // 紧急程度，1-5，数值越大紧急程度越高
	UndertakerID string  `gorm:"type:varchar(100)"`  // 承接人，接单的管理员id，不允许为空（默认就写个null)
	IfRubbish    int     `json:"if_rubbish"`         // 状态，默认1,计入评价状况
	IsAnonymous  bool    `json:"is_anonymous"`       // 是否匿名，默认false（不匿）
	Status       bool    `json:"status"`             // 是否已经处理，默认false（未
	Grade        int     `json:"grade"`              // 评分
	GradeContent string  `json:"grade_content"`      // 评价内容
	Student      Student `gorm:"foreignKey:UserID;references:UserID"`
	Admin        Admin   `gorm:"foreignKey:UndertakerID;references:UserID"`
}
