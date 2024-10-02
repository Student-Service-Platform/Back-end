package controllers

import (
	"Back-end/services"
	"Back-end/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthAPI 负责处理用户注册、登录和信息修改后的操作

// 注册功能等等再写，邮件那个东西要整的有点多

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
				utils.JsonResponse(ctx, 200, 200200, "登录成功", gin.H{
					"username": user.UserId,
					"type":     user.Type,
					"token":    utils.GenerateJWT(user.UserId, user.Type),
				})
			}
		}
	}
}
