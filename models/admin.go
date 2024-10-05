package models

import "gorm.io/gorm"

// Admin 管理员
type Admin struct {
	gorm.Model         // 自行管理的created_at, updated_at, deleted_at
	UserID     string  `json:"user_id" gorm:"type:varchar(100);uniqueIndex"`
	Username   string  `json:"username"`
	Password   string  `json:"-"`
	Phone      string  `json:"phone"`
	Mail       string  `json:"mail"`
	Type       int     `json:"type"`
	Avatar     string  `json:"avatar"`
	IfDel      bool    `json:"if_del"`
	HadDone    uint    `json:"had_done"`
	Evalutaion float32 `json:"evalutaion"`
}
