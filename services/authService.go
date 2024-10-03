package services

// 这个文件处理：注册，登录
import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"

	"golang.org/x/crypto/bcrypt"
)

// 检测用户是否存在的函数
type UserInfo struct {
	UserID string
}

func CheckUserExistByUserID(UserID string, table string) error {
	result := database.DB.Table(table).Where("user_id = ?", UserID).First(&UserInfo{})
	utils.LogError(result.Error)
	return result.Error
}

// 添加学生用户
var NewStudent models.Student

func AddStudent(UserID string, Username string, RawPassword string) error {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RawPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(err)
		return err
	}

	NewStudent = models.Student{
		UserID:   UserID,
		Username: Username,
		Password: string(hashedPassword),
		Phone:    "",
		Mail:     "",
		IfDel:    false,
		Avatar:   "https://imgse.com/i/pA89yOe",
		Type:     1,
	}

	result := database.DB.Table("students").Create(&NewStudent)
	return result.Error
}

// 获取用户的函数
type LoginInfo struct {
	UserID   string
	Password string
	Type     int
}

func GetUserByUserID(UserID string, table string) (*LoginInfo, error) {
	var userinfo LoginInfo
	result := database.DB.Table(table).Where("user_id = ?", UserID).First(&userinfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userinfo, nil
}

// 比对密码
func CheckPassword(pwdinput string, pwddb string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwddb), []byte(pwdinput))
	utils.LogError(err)
	return err == nil
}
