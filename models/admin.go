package models // Admin
import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	// 头像，apifox的数据库模型里没有图片，暂时采用url文本代替
	Avatar *string `json:"avatar,omitempty"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 邮箱，联系方式1
	Email string `json:"email"`
	// 手机号，联系方式2（选填，预留接口）
	Number string `json:"number"`
	// 用户类型，（1普通用户）2管理员3超级管理员
	Type int64 `json:"type"`
	// 修改时间
	UpdatedAt string `json:"updated_at"`
	// 管理员id，A+xxxx
	UserID string `json:"user_id"`
	// 姓名，用户名
	Username string `json:"username"`
}
