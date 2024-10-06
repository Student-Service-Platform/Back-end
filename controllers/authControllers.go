package controllers

import (
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

	// 检查学号是否合法（纯数字）
	getint, err := strconv.Atoi(data.UserID)
	if err != nil {
		utils.JsonResponse(ctx, 200, 200501, "你这学号有问题啊", nil)
		utils.LogError(err)
	} else if getint <= 0 {
		utils.JsonResponse(ctx, 200, 200501, "你这学号有问题啊", nil)
	} else { //确定学号没问题之后检验密码
		if len(data.Password) < 8 || len(data.Password) > 16 {
			utils.JsonResponse(ctx, 200, 200502, "你这密码有问题啊", nil)
		} else {
			//检验是否存在
			err = services.CheckUserExistByUserID(data.UserID, "students")
			if err != nil && err != gorm.ErrRecordNotFound { //如果发生了错误
				utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
			} else if err == nil { //如果找到了
				utils.JsonResponse(ctx, 200, 200504, "你这学号已经被注册了", nil)
			} else {
				//补齐struct然后塞进去，这里只负责传入相关信息
				err = services.AddStudent(data.UserID, data.Username, data.Password)
				if err != nil {
					utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
					utils.LogError(err)
				} else {
					utils.JsonResponse(ctx, 200, 200200, "注册成功", nil)
				}
			}
		}
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
				ctx.SetCookie("user_id", data.UserID, 3600000, "/", "localhost", false, false)
				utils.JsonResponse(ctx, 200, 200200, "登录成功", gin.H{
					"user_id": user.UserID,
					"type":    user.Type,
					"token":   utils.GenerateJWT(user.UserID, user.Type),
				})
			}
		}
	}
}
