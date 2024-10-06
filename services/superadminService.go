package services

// 这个文件处理：注册，登录
import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 添加管理员用户
var NewAdmin models.Admin

func AddAdmin(UserID string, Username string, RawPassword string) error {
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RawPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(err)
		return err
	}

	NewAdmin = models.Admin{
		UserID:     UserID,
		Username:   Username,
		Password:   string(hashedPassword),
		Phone:      "",
		Mail:       "",
		IfDel:      false,
		Avatar:     "https://imgse.com/i/pA89yOe",
		Type:       2,
		HadDone:    0,
		Evalutaion: 0,
	}

	result := database.DB.Table("admins").Create(&NewAdmin)
	return result.Error
}

// 删除用户
func DelUser(userID string, table string) error {
	var count1 int64
	var count2 int64

	// 统计 student 表中 user_id 的记录数
	if err := database.DB.Model(&models.Student{}).Where("user_id = ?", userID).Count(&count1).Error; err != nil {
		return err
	}

	// 统计 admin 表中 user_id 的记录数
	if err := database.DB.Model(&models.Admin{}).Where("user_id = ?", userID).Count(&count2).Error; err != nil {
		return err
	}

	// 如果两个表中都没有该用户，则返回错误信息
	if count1+count2 == 0 {
		return errors.New("该用户不存在")
	}

	// 如果在 student 表中有记录，则删除 student 记录
	if count1 != 0 {
		if err := database.DB.Where("user_id = ?", userID).Delete(&models.Student{}).Error; err != nil {
			return err
		}
		return nil
	}

	// 否则删除 admin 记录
	if err := database.DB.Where("user_id = ?", userID).Delete(&models.Admin{}).Error; err != nil {
		return err
	}

	return nil
}

type getusers struct {
	UserID   string
	UserName string
	Password string
	Mail     string
	Phone    string
	IfDel    bool
}

func GetAll(table string, offset int, limit int) ([]getusers, error) {

	if table == "admins" {
		var Users []models.Admin
		if err := database.DB.Offset(offset).Limit(limit).
			Find(&Users).Error; err != nil {
			return nil, err
		}
		var UserInfos []getusers
		for _, req := range Users {
			userInfo := getusers{
				UserID:   req.UserID,
				UserName: req.Username,
				Password: req.Password,
				Mail:     req.Mail,
				Phone:    req.Phone,
				IfDel:    req.IfDel,
			}
			UserInfos = append(UserInfos, userInfo)
		}
		return UserInfos, nil
	} else {
		var Users []models.Student
		if err := database.DB.Offset(offset).Limit(limit).
			Find(&Users).Error; err != nil {
			return nil, err
		}
		var UserInfos []getusers
		for _, req := range Users {
			userInfo := getusers{
				UserID:   req.UserID,
				UserName: req.Username,
				Password: req.Password,
				Mail:     req.Mail,
				Phone:    req.Phone,
				IfDel:    req.IfDel,
			}
			UserInfos = append(UserInfos, userInfo)
		}
		return UserInfos, nil
	}
}

//Format("2006-01-02 15:04"),
