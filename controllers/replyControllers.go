package controllers

import (
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 管理员回复反馈
func ReplyRequest(ctx *gin.Context) {
	currentFeedbackID := ctx.Param("id") //这得是int啊
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}

	intFeedbackID, err := strconv.Atoi(currentFeedbackID)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "客户端报告了数据型错误", nil)
		return
	}

	type replyContent struct {
		Content string `json:"content"`
	}
	var content replyContent
	if err := ctx.BindJSON(&content); err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "客户端报告了数据型错误", nil)
		return
	}

	// TODO: 添加回复反馈的逻辑

	if 2 != userType && 3 != userType {
		utils.JsonResponse(ctx, 200, 200401, "不要说我们没警告过你……你的权限不对", nil)
	} else {
		err = services.CheckRequestReplyExistByID(int64(intFeedbackID))
		if err != nil && err != gorm.ErrRecordNotFound { //如果发生了错误
			utils.JsonResponse(ctx, 200, 200503, "咱整了点问题，晚点再试吧", nil)
		} else if err == nil { //如果找到了
			utils.JsonResponse(ctx, 200, 200510, "有回复在你之前赶到了！", nil)
		} else {
			//补齐struct然后塞进去，这里只负责传入相关信息
			err = services.CreateRequestReply(content.Content, currentUserID, int64(intFeedbackID))
			if err != nil {
				utils.JsonResponse(ctx, 200, 200509, "出了点点问问题题题题", nil)
				utils.LogError(err)
			} else {
				utils.JsonResponse(ctx, 200, 200200, "好东西就要来了！刷新查看（大概吧）", nil)
			}
		}
	}
}
