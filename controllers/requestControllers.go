package controllers

import (
	"Back-end/models"
	"Back-end/services"
	"Back-end/utils"

	"github.com/gin-gonic/gin"
)

// postController，对于Requsest的处理
// 在 userController中已经有了获取用户信息的函数parseContext

// 解析中间件传递的上下文信息

// Student账户部分----------
// 创建Request
type createRequest struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Category    int64  `json:"category"`
	Description string `json:"description"`
	IsUrgent    int64  `json:"is_urgent"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func CreateRequest(ctx *gin.Context) {
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}
	var request createRequest
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}

	switch userType {
	case 1: // 普通账户
		err = services.CreateRequest(models.Request{
			UserID:       currentUserID,
			Title:        request.Title,
			Description:  request.Description,
			Category:     request.Category,
			Urgency:      request.IsUrgent,
			UndertakerID: "",
			IsAnonymous:  request.IsAnonymous,
			IfRubbish:    1,
			Status:       false,
			Grade:        0,
			GradeContent: "",
		})
		if err != nil {
			utils.JsonResponse(ctx, 200, 200504, "服务器出错，我们都有不顺利的时候，尝试在晚点", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "创建成功", nil)
		}
	case 2, 3: //管理员账户
		utils.JsonResponse(ctx, 200, 200401, "客户端报告您可能不是普通人，换个账户试试", nil)
	default:
		utils.JsonResponse(ctx, 200, 200506, "你可能没有合适的权限，坐和放宽。", nil)
	}
}
