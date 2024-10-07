package controllers

import (
	"Back-end/services"
	"Back-end/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

// userController 用户控制器，用户对于自己信息的增删改查

// 解析中间件传递的上下文信息
func parseContext(ctx *gin.Context) (string, int, string, error) {
	currentUserID := ctx.GetString("userID")
	userType := ctx.GetInt("type")
	var table string

	switch userType {
	case 1:
		table = "students"
		break
	case 2, 3:
		table = "admins"
		break
	default:
		utils.JsonResponse(ctx, 200, 401, "需要滚回到以前的用户类型", nil)
		return "", 0, "", fmt.Errorf("用户类型无效")
	}
	// 获取用户信息直接 currentUserID, userType, table, err := parseContext(ctx)
	return currentUserID, userType, table, nil
}

// 获取用户信息（用户获取用户，管理员获取管理员，不加参数获取自己的）
func GetProfile(ctx *gin.Context) {
	targetUserID := ctx.Query("user_id")

	currentUserID, userType, table, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}

	if targetUserID != currentUserID {
		utils.JsonResponse(ctx, 403, 200506, "你可能没有合适的权限，坐和放宽。", nil)
		return
	}

	fmt.Print(targetUserID)

	var result services.ReturnUserInfo

	result, err = services.GetProfileByID(targetUserID, table)
	if err != nil {
		utils.JsonResponse(ctx, 403, 200506, "你可能没有合适的权限，坐和放宽。", nil)
	} else {
		utils.JsonResponse(ctx, 200, 200200, "success", gin.H{
			"user_id": result.UserID,
			"name":    result.Username,
			"phone":   result.Phone,
			"mail":    result.Mail,
			"avatar":  result.Avatar,
			"type":    userType,
		})
	}
}

// 修改个人信息
type upDateInfo struct {
	Object   string `json:"object"`
	NewValue string `json:"new_value"`
}

func UpdateProfile(ctx *gin.Context) {
	currentUserID, _, table, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}

	var data upDateInfo
	// 解析参数并且绑定到 struct 中
	err = ctx.ShouldBindJSON(&data)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}
	fmt.Println(table)
	switch data.Object {
	case "username", "password", "mail", "phone", "avatar":

		switch data.Object {
		case "username":
			if data.NewValue == "" || data.NewValue == "匿名用户" {
				utils.JsonResponse(ctx, 200, 200506, "你小子，用户名有问题", nil)
				return
			}
			break
		case "password":
			if len(data.NewValue) < 8 || len(data.NewValue) > 16 {
				utils.JsonResponse(ctx, 200, 200502, "你小子，密码长度有问题啊", nil)
				return
			}
			break
		case "mail":
			if !utils.IsValidMail(data.NewValue) {
				utils.JsonResponse(ctx, 200, 200506, "你小子，邮箱有问题啊", nil)
				return
			}
			break
		case "phone":
			if !utils.IsValidPhone(data.NewValue) {
				utils.JsonResponse(ctx, 200, 200506, "你小子，手机号有问题啊", nil)
				return
			}
			break
		case "avatar":
			if !utils.IsValidURL(data.NewValue) {
				utils.JsonResponse(ctx, 200, 200506, "你小子，连个图片链接都传不明白？", nil)
				return
			}
		default:
			utils.JsonResponse(ctx, 200, 200508, "让我们重回正轨！请选择有效的修改字段。", nil)
		}

		err = services.UpdateProfile(currentUserID, table, data.Object, data.NewValue)
		if err != nil {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200506, "这下尴尬了。我们好好像出了点问题。", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "修改成功！大概得刷新页面生效吧。", nil)
		}
		break
	default:
		utils.JsonResponse(ctx, 200, 200508, "让我们重回正轨！请选择有效的修改字段。", nil)
	}
}
