package controllers

import (
	"Back-end/models"
	"Back-end/services"
	"Back-end/utils"
	"fmt"
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

// 接单函数
// 处理请求
func HandleRequst(ctx *gin.Context) {
	// 获取当前反馈ID
	currentFeedbackID := ctx.Param("id")
	// 获取操作类型
	action := ctx.Query("action")
	// 将当前反馈ID转换为整数
	intFeedbackID, err := strconv.Atoi(currentFeedbackID)
	if err != nil {
		// 记录错误日志
		utils.LogError(err)
		// 返回错误信息
		utils.JsonResponse(ctx, 200, 200503, "报告数数数据据据型错误", nil)
		return
	}
	// 将操作类型转换为整数
	intaction, err := strconv.Atoi(action)
	if err != nil {
		// 记录错误日志
		utils.LogError(err)
		// 返回错误信息
		utils.JsonResponse(ctx, 200, 200503, "报告数数数据据据型错误", nil)
		return
	}

	// 解析上下文，获取当前用户ID和用户类型
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		// 记录错误日志
		utils.LogError(err)
		return
	}

	// 判断用户类型是否为2或3
	if userType != 2 && userType != 3 {
		// 返回错误信息
		utils.JsonResponse(ctx, 200, 200401, "头抬起，你的权限不对", nil)
		return
	}

	// 判断当前反馈ID是否已被处理
	existUserID, err := services.IsHandled(intFeedbackID)
	if err != nil {
		// 记录错误日志
		utils.LogError(err)
		// 返回错误信息
		utils.JsonResponse(ctx, 200, 200503, "没关系，我们都有不顺利的时候", nil)
		return
	}

	// 判断当前用户是否已处理过该反馈ID
	if existUserID != "" && existUserID != currentUserID {
		// 返回错误信息
		utils.JsonResponse(ctx, 200, 200511, "有人在你之前遇见了！", nil)
		return
	}

	// 根据操作类型进行不同的处理
	switch intaction {
	case 1:
		// 处理请求
		err1 := services.HandleRequest(intFeedbackID, currentUserID)
		var err2 error
		// 如果当前反馈ID未被处理过，更新管理员已处理数量
		if existUserID == "" {
			err2 = services.UpdateAdminHaddone(currentUserID, 1)
		}
		// 如果处理过程中出现错误，记录错误日志并返回错误信息
		if err1 != nil || err2 != nil {
			utils.LogError(err1)
			utils.LogError(err2)
			utils.JsonResponse(ctx, 200, 200504, "数据库出现在问题，尝试在晚点", nil)
		}
	case 0:
		// 取消处理请求
		err1 := services.HandleRequest(intFeedbackID, "")
		err2 := services.UpdateAdminHaddone(currentUserID, -1)
		// 如果处理过程中出现错误，返回错误信息
		if err1 != nil || err2 != nil {
			utils.JsonResponse(ctx, 200, 200504, "数据库出现在问题，尝试在晚点", nil)
		}
	default:
		// 记录错误日志
		utils.LogError(fmt.Errorf("坏东西来了：未指定的接单操作"))
	}
}

// 评价处理结果（是这个值是写在request里头的，以及要同步更新管理员那里的总分）
// 处理评分请求
func GradeHandleRequest(ctx *gin.Context) {
	// 获取请求参数id
	id := ctx.Param("id")
	// 将id转换为int类型
	intID, err := strconv.Atoi(id)
	// 如果转换失败，记录错误并返回
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "报告反馈数数据据据型错误", nil)
		return
	}

	// 定义输入结构体
	var input struct {
		Grade        int    `json:"grade"`
		GradeContent string `json:"grade_content"`
	}
	// 绑定json数据到输入结构体
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "数据错误在传入", nil)
		return
	}

	// 根据id获取请求
	currentRequest, err := services.GetRequestByID(intID)
	// 如果获取失败，记录错误并返回
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200512, "很久很久以前，大概有这么一个反馈", nil)
		return
	}

	// 解析上下文，获取当前用户id和用户类型
	currentUserID, userType, _, err := parseContext(ctx)
	// 如果解析失败，记录错误并返回
	if err != nil {
		utils.LogError(err)
		return
	}

	// 如果用户类型为1或者请求的执行者id不为空且请求不是垃圾请求且请求的用户id等于当前用户id
	if userType != 1 && currentRequest.UndertakerID != "" && currentRequest.IfRubbish == 0 && currentRequest.UserID == currentUserID {
		// 更新请求的评分和评分内容
		currentRequest.Grade = input.Grade
		currentRequest.GradeContent = input.GradeContent
		// 更新请求的评分
		if err := services.UpdateRequestEvaluation(&currentRequest); err != nil {
			utils.JsonResponse(ctx, 200, 200403, "评价失败，请检查身份和请求状态", nil)
			return
		}
		// 返回成功信息
		utils.JsonResponse(ctx, 200, 200200, "评价成功", nil)
	} else {
		// 返回失败信息
		utils.JsonResponse(ctx, 200, 200403, "评价失败，请检查身份或请求状态", nil)
		return
	}

	// 更新执行者的评分
	if err := services.UpdateAdminEvaluation(currentRequest.UndertakerID, input.Grade); err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200513, "数据库出现在问题", nil)
		return
	}
}

// 将request标记为垃圾
func MarkRequest(ctx *gin.Context) {
	// 获取请求参数id
	id := ctx.Param("id")
	// 将id转换为int类型
	intID, err := strconv.Atoi(id)
	// 如果转换失败，记录错误并返回
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "数据错误在传入", nil)
	}
	// 解析上下文，获取当前用户id和用户类型
	_, userType, _, err := parseContext(ctx)
	// 如果解析失败，记录错误并返回
	if err != nil {
		utils.LogError(err)
		return
	}
	switch userType {
	case 2: //普通管理员内
		// 标记请求为垃圾
		err = services.MarkRequest(intID)
		if err != nil {
			utils.LogError(err)
			utils.JsonResponse(ctx, 200, 200513, "坐和放宽，数据库出现了一些问题", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "标记成功", nil)
		}
		break
	case 3: //超级管理员
		confirm := ctx.Query("confirmation")
		if "true" == confirm { //超管确认垃圾
			err = services.MarkRequest(intID)
			if err != nil {
				utils.LogError(err)
				utils.JsonResponse(ctx, 200, 200513, "坐和放宽，数据库又整了点错误出来", nil)
			}
			return
		} else if "false" == confirm {
			err = services.StatueRequest(intID)
			if err != nil {
				utils.LogError(err)
				utils.JsonResponse(ctx, 200, 200513, "坐和放宽，数据库又整了点错误出来", nil)
			}
			return
		}
		break
	default:
		utils.JsonResponse(ctx, 200, 200403, "权限不足", nil)
	}
}
