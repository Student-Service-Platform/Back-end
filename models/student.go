package models

import "gorm.io/gorm"

// Student 学生
type Student struct {
	gorm.Model
	UserID   string `json:"user_id" gorm:"index"` // 学号，学号
	Username string `json:"username"`         // 姓名，用户名(显示名称)
	Password string `json:"-"`                // 密码
	Phone    string `json:"phone,omitempty"`  // 手机号，联系方式2（选填，预留接口）
	Mail     string `json:"mail"`             // 邮箱，联系方式一
	Type     int    `json:"type"`             // 直接写死=1，有这个数值是为了和后面管理员那边统一返回数据，目前未发现影响其余函数
	Avatar   string `json:"avatar,omitempty"` // 头像，apifox的数据库模型里没有图片，暂时采用url文本代替
	IfDel    bool   `json:"if_del"`           // `被删掉了嘛
}
