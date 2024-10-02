package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model         //自行管理的created_at, updated_at, deleted_at
	UserID     string  `json:"user_id"`
	Name       string  `json:"name"`
	Password   string  `json:"-"`
	Phone      string  `json:"phone"`
	Mail       string  `json:"mail"`
	Type       uint    `json:"type"`
	Avatar     string  `json:"avatar"`
	IfDel      bool    `json:"if_del"`
	HadDone    uint    `json:"haddone"`
	Evalutaion float32 `json:"evalutaion"`
}
