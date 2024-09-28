package models

import "gorm.io/gorm"

// Student 学生
type Student struct {
	gorm.Model

	// 头像，apifox的数据库模型里没有图片，暂时采用url文本代替
	Avatar *string `json:"avatar,omitempty"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 邮箱，联系方式一
	Mail string `json:"mail"`
	// 密码
	Password string `json:"password"`
	// 手机号，联系方式2（选填，预留接口）
	Phone *string `json:"phone,omitempty"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
	// 学号，学号
	UserID string `json:"user_id"`
	// 姓名，用户名
	Username string `json:"username"`
}
