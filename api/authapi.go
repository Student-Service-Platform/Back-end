package api

import (
	"ServerPlatform/utils"

	"github.com/gin-gonic/gin"
)

// AuthAPI 负责处理用户注册、登录和信息修改后的操作

// 注册功能等等再写，邮件那个东西要整的有点多

type loginData struct {
	username string `json:"username" binding:"required"`
	password string `json:"password" binding:"required"`
	isadmin  bool   `json:"isadmin" binding:"required"`
}

func Login(ctx *gin.Context) {
	var data loginData
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		utils.LogError(err)
	}
}
