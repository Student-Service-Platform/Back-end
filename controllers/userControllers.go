package controllers

import (
	"Back-end/services"
	"Back-end/utils"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(ctx *gin.Context) {
	currentUserID := ctx.GetString("user_id")
	userType := ctx.GetInt("type")
	targetUserID := ctx.Query("user_id")
	if targetUserID == "" {
		targetUserID = currentUserID
	}

	var result services.ReturnUserInfo
	var err error

	switch userType {
	case 1:
		table := "students"
		result, err = services.GetProfileByID(targetUserID, table)
		if err != nil {
			utils.JsonResponse(ctx, 200, 200506, "这真是让人尴尬，请坐和放宽。", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "成功获取用户信息", gin.H{
				"user_id": result.UserID,
				"name":    result.Username,
				"phone":   result.Phone,
				"mail":    result.Mail,
				"avatar":  result.Avatar,
				"type":    userType,
			})
		}
		break
	case 2, 3:
		table := "admins"
		result, err = services.GetProfileByID(targetUserID, table)
		if err != nil {
			utils.JsonResponse(ctx, 200, 200506, "这真是让人尴尬，请坐和放宽。", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "成功获取用户信息", gin.H{
				"user_id": result.UserID,
				"name":    result.Username,
				"phone":   result.Phone,
				"mail":    result.Mail,
				"avatar":  result.Avatar,
				"type":    userType,
			})
		}
		break
	default:
		utils.JsonResponse(ctx, 200, 401, "未授权的访问", nil)
		break
	}
}
