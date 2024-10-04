package controllers

import (
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
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
		existUserID, err := services.IsHandled(int64(intFeedbackID))
		if nil != err {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200503, "没关系，我们都有不顺利的时候", nil)
		} else if "" != existUserID || currentUserID != existUserID {
			utils.JsonResponse(ctx, 200, 200510, "有回复在你之前赶到了！", nil)
		} else {
			//补齐struct然后塞进去，这里只负责传入相关信息
			err1, err2 := services.CreateRequestReply(content.Content, currentUserID, int64(intFeedbackID)), services.HandleRequest(int64(intFeedbackID), currentUserID)
			if nil != err1 || nil != err2 {
				utils.JsonResponse(ctx, 200, 200509, "出了点点问问题题题题", nil)
				utils.LogError(err1)
				utils.LogError(err2)
			} else {
				utils.JsonResponse(ctx, 200, 200200, "好东西就要来了！刷新查看（大概吧）", nil)
			}
		}
	}
}
