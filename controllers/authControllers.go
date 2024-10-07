package controllers

import (
	"Back-end/config"
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthAPI 负责处理用户注册、登录的操作

// 注册功能等等再完善，邮件那个东西要整的有点多，我这里只写了基本的学号+用户名+密码的模式，但是测试应该够用了
type registerData struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	MailAuth bool   `json:"mail_auth"` //还没弄完
	Phone    string `json:"phone"`
	Mail     string `json:"mail" binding:"required"`
}

func Register(ctx *gin.Context) {
	var data registerData
	// 解析参数并且绑定到 struct 中
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}

	// 将字符串UserID转换为整数
	userIDInt, err := strconv.Atoi(data.UserID)
	// 如果转换失败或者UserID小于等于0，则返回错误
	if err != nil || userIDInt <= 0 {
		utils.JsonResponse(ctx, 200, 200501, "你这学号有问题啊", err)
		utils.LogError(err)
		return
	}

	// 如果密码长度不在8到16位之间，则返回错误
	if len(data.Password) < 8 || len(data.Password) > 16 {
		utils.JsonResponse(ctx, 200, 200502, "你这密码长度有问题啊", nil)
		return
	}

	// 根据UserID检查用户是否存在
	err = services.CheckUserExistByUserID(data.UserID, "students")
	// 如果存在，则返回错误
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:

			if data.Username == "" || data.Username == "匿名用户" {
				utils.JsonResponse(ctx, 200, 200506, "你小子，用户名有问题", nil)
				return
			}
			// 如果邮箱格式不正确，则返回错误
			if !utils.IsValidMail(data.Mail) {
				utils.JsonResponse(ctx, 200, 200503, "你这邮箱有问题啊", nil)
				return
			}

			// 如果手机号格式不正确，则返回错误
			if !utils.IsValidPhone(data.Phone) {
				utils.JsonResponse(ctx, 200, 200503, "你这手机号有问题啊", nil)
				return
			}

			// 添加学生信息
			err = services.AddStudent(data.UserID, data.Username, data.Password, data.Phone, data.Mail)
			// 如果添加失败，则返回错误
			if err != nil {
				utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", err)
				utils.LogError(err)
				return
			}

			// 返回成功
			utils.JsonResponse(ctx, 200, 200200, "注册成功", nil)
		default:
			utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", err)
			utils.LogError(err)
			return
		}
	} else {
		// 如果学生信息已存在，则返回错误
		utils.JsonResponse(ctx, 200, 200504, "你这学号已经被注册了", nil)
		return
	}
}

// 登录功能↓
type loginData struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

func Login(ctx *gin.Context) {
	var data loginData
	// 解析参数并且绑定到 struct 中
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}

	// 管理员/普通用户登录
	var table string
	if data.IsAdmin {
		table = "admins"
	} else {
		table = "students"
	}

	err = services.CheckUserExistByUserID(data.UserID, table)
	if err != nil { //如果发生了错误
		if err == gorm.ErrRecordNotFound { //如果是未找到
			utils.JsonResponse(ctx, 200, 200505, "你这学号有问题啊", nil)
		} else { //如果找到了还是发生了错误
			utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
			utils.LogError(err)
			return
		}
	} else { //用户存在，检测过程中没有出现错误
		//确认用户存在后再检测密码
		user, err := services.GetUserByUserID(data.UserID, table)

		if err != nil {
			utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
			utils.LogError(err)
			return
		} else {
			flag := services.CheckPassword(data.Password, user.Password) //调用service层中的检测密码函数
			if !flag {
				utils.JsonResponse(ctx, 200, 200504, "你这密码有问题啊", nil)
			} else {
				ctx.SetCookie(
					"user_id", data.UserID,
					config.Config.GetInt("cookies.maxAge"),
					config.Config.GetString("cookies.path"),
					config.Config.GetString("cookies.domain"),
					config.Config.GetBool("cookies.secure"),
					config.Config.GetBool("cookies.httpOnly"))
				ctx.Set("userID", user.UserID)
				ctx.Set("type", user.Type)
				utils.JsonResponse(ctx, 200, 200200, "登录成功", gin.H{
					"user_id": user.UserID,
					"type":    user.Type,
					"token":   utils.GenerateJWT(user.UserID, user.Type),
				})
			}
		}
	}
}
