package services

// 这个文件处理：注册，登录
import (
	"Back-end/database"
	"Back-end/utils"

	"golang.org/x/crypto/bcrypt"
)

// 检测用户是否存在的函数
type UserInfo struct {
	UserID string
}

func CheckUserExistByUserID(UserId string, table string) error {
	result := database.DB.Table(table).Where("user_id = ?", UserId).First(&UserInfo{})
	utils.LogError(result.Error)
	return result.Error
}

// 获取用户的函数
type LoginInfo struct {
	UserId   string
	Password string
	Type     int
}

func GetUserByUserID(UserId string, table string) (*LoginInfo, error) {
	var userinfo LoginInfo
	result := database.DB.Table(table).Where("user_id = ?", UserId).First(&LoginInfo{})
	if result.Error != nil {
		return nil, result.Error
	}
	return &userinfo, nil
}

// 比对密码
func CheckPassword(pwdinput string, pwddb string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwddb), []byte(pwdinput))
	return err == nil
}
