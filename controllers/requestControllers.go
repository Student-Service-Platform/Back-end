package controllers

import (
	"Back-end/models"
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// postController，对于Requsest的处理
// 在 userController中已经有了获取用户信息的函数parseContext

// 解析中间件传递的上下文信息

// 创建Request：要求userType=1
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

// 查看Request 不需要登录验证
func GetRequest(ctx *gin.Context) {
	UserID := ctx.Query("user_id")
	pageStr := ctx.Query("page")
	perPageStr := ctx.Query("limit")

	if "" == pageStr {
		pageStr = "1"
	}

	if "" == perPageStr {
		perPageStr = "15"
	}
	// ... 页面参数转换 ...
	page, err1 := strconv.Atoi(pageStr)       // 将 page 字符串转换为整数
	perPage, err2 := strconv.Atoi(perPageStr) // 将 per_page 字符串转换为整数
	if err1 != nil || err2 != nil {
		// 处理错误
		utils.LogError(err1)
		utils.LogError(err2)
		utils.JsonResponse(ctx, 200, 200503, "这下尴尬了。。。我们正在让一切重回正轨", nil)
	}
	if page <= 0 || perPage <= 0 {
		page, perPage = 1, 15 //默认设置
	}
	offset := (page - 1) * perPage // 计算偏移量

	if UserID == "" { //可以看得到所有的Request，包括匿名的
		requests, err := services.GetAllRequest(offset, perPage)
		if err != nil {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200504, "服务器出错，我们都有不顺利的时候，尝试在晚点", nil)
		} else {
			if len(requests) == 0 {
				utils.JsonResponse(ctx, 200, 200200, "还没有发过哦", nil)
			} else {
				utils.JsonResponse(ctx, 200, 200200, "success", requests)
			}
		}
	} else {
		//看特定用户的Request，不包括匿名的
		requests, err := services.GetAllRequest(offset, perPage)
		if err != nil {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200504, "服务器出错，我们都有不顺利的时候，尝试在晚点", nil)
		} else {
			if len(requests) == 0 {
				utils.JsonResponse(ctx, 200, 200200, "还没有发过哦", nil)
			} else {
				utils.JsonResponse(ctx, 200, 200200, "success", requests)
			}
		}
	}
}

//
