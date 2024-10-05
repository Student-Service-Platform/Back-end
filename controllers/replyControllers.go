package controllers

import (
	"Back-end/models"
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 回复反馈
func ReplyRequest(ctx *gin.Context) {
	currentFeedbackID := ctx.Param("id")
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "客户端报告数数据型错误", nil)
		return
	}

	// 将字符串类型的反馈ID转换为int类型
	intFeedbackID, err := strconv.Atoi(currentFeedbackID)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "客户端报告了数据据型错误", nil)
		return
	}

	// 定义回复内容结构体
	type replyContent struct {
		Content string `json:"content"`
	}
	var content replyContent
	// 解析请求中的JSON数据
	if err := ctx.BindJSON(&content); err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "客户端报告了数据型型错误", nil)
		return
	}

	// 判断用户类型
	switch userType {
	case 1, 2, 3:
		reply := models.Reply{
			Father:     0,
			RequestID:  int64(intFeedbackID),
			Content:    content.Content,
			Respondent: currentUserID,
		}
		err := services.CreateReply(&reply)
		if err != nil {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200509, "出了点点问题", nil)
			return
		}
		utils.JsonResponse(ctx, 200, 200200, "回复成功", nil)
		break
	default:
		utils.JsonResponse(ctx, 200, 200401, "不要说我们没警告过你……你的权限不对", nil)
		break
	}
}
