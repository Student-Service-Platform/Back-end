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
	Category    int    `json:"category"`
	Description string `json:"description"`
	IsUrgent    bool   `json:"is_urgent"`
	IsAnonymous bool   `json:"is_anonymous"`
}

func CreateRequest(ctx *gin.Context) {
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}
	var request createRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}
	switch userType {
	case 1: // 普通账户
		if request.Title == "" || request.Description == "" {
			utils.JsonResponse(ctx, 200, 200503, "标题和描述不能为空", nil)
			return
		}
		var urgent int
		if request.IsUrgent {
			urgent = 1
		} else {
			urgent = 0
		}
		err = services.CreateRequest(models.Request{
			UserID:       currentUserID,
			Title:        request.Title,
			Description:  request.Description,
			Category:     request.Category,
			Urgency:      urgent,
			UndertakerID: "null",
			IsAnonymous:  request.IsAnonymous,
			IfRubbish:    1,
			Status:       true,
			Grade:        0,
			GradeContent: "",
		})
		if err != nil {
			utils.JsonResponse(ctx, 200, 200504, "服务器出错，我们都有不顺利的时候，尝试在晚点", nil)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "创建成功", nil)
		}
		break
	case 2, 3: // 管理员账户
		utils.JsonResponse(ctx, 200, 200401, "客户端报告您可能不是普通人，换个账户试试", nil)
		break
	default:
		utils.JsonResponse(ctx, 200, 200506, "你可能没有合适的权限，坐和放宽。", nil)
		break
	}
}

// 查看Request 不需要登录验证

func GetAllRequests(ctx *gin.Context) {
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
	if UserID == "" {              //可以看得到所有的Request，包括匿名的
		requests, err := services.GetAllRequests(offset, perPage)
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
		requests, err := services.GetRequestsByUserID(UserID, offset, perPage)
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

// 获取选取过的反馈列表
func GetSelectedFeedback(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	perPageStr := ctx.Query("limit")
	status := ctx.Query("status")   // 0未处理 >=1已处理
	rubbish := ctx.Query("rubbish") // 0垃圾 >=1还不是垃圾

	//默认值设置 pageStr=1，perPageStr=15，status=0（未处理），rubbish=1（不是垃圾）
	if "" == pageStr {
		pageStr = "1"
	}

	if "" == perPageStr {
		perPageStr = "15"
	}

	if "" == status {
		status = "0"
	}

	if "" == rubbish {
		rubbish = "1"
	}

	// ... 页面参数转换 ...
	page, err1 := strconv.Atoi(pageStr)       // 将 page 字符串转换为整数
	perPage, err2 := strconv.Atoi(perPageStr) // 将 per_page 字符串转换为整数
	state, err3 := strconv.Atoi(status)
	if_rubbish, err4 := strconv.Atoi(rubbish)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		// 处理错误
		utils.LogError(err1)
		utils.LogError(err2)
		utils.LogError(err3)
		utils.LogError(err4)
		utils.JsonResponse(ctx, 200, 200515, "错误在输入。。。我们正在让一切重回正轨", nil)
	} // 这一串错误处理有没有更优雅的方法？
	if page <= 0 {
		page = 1 //数值不对的情况下默认设置
	}

	if perPage <= 0 {
		perPage = 15
	}

	if state != 0 {
		state = 1
	}

	if state != 0 { //
		state = 1
	} //已处理的

	offset := (page - 1) * perPage // 计算偏移量
	requests, err := services.GetSelectRequests(offset, perPage, if_rubbish, state)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200504, "服务器出错，我们都有不顺利的时候，尝试在晚点", nil)
		return
	}
	if len(requests) == 0 {
		utils.JsonResponse(ctx, 200, 200200, "还没有发过哦", nil)
	} else {
		utils.JsonResponse(ctx, 200, 200200, "success", requests)
	}
}

// 接单函数
// 处理请求
func HandleRequest(ctx *gin.Context) {
	// 获取当前反馈ID
	currentFeedbackID := ctx.Param("id")
	// 获取操作类型
	action := ctx.Query("action")

	// 将当前反馈ID转换为整数
	intFeedbackID, err := strconv.Atoi(currentFeedbackID)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "报告数数据据据型错误", nil)
		return
	}

	// 将操作类型转换为整数
	intAction, err := strconv.Atoi(action)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "报告数数数据据型错误", nil)
		return
	}

	// 解析上下文，获取当前用户ID和用户类型
	currentUserID, userType, _, err := parseContext(ctx)
	if err != nil {
		utils.LogError(err)
		return
	}

	// 判断用户类型是否为2或3
	if userType != 2 && userType != 3 {
		utils.JsonResponse(ctx, 200, 200401, "头抬起，你的权限不对", nil)
		return
	}

	// 判断当前反馈ID是否已被处理
	existUserID, err := services.IsHandled(intFeedbackID)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "没关系，我们都有不顺利的时候", nil)
		return
	}

	// 判断是不是有人已经处理过这个反馈了
	if existUserID != "null" && existUserID != currentUserID {
		utils.JsonResponse(ctx, 200, 200511, "有人在你之前遇见了！", nil)
		return
	}

	// 根据操作类型进行不同的处理
	switch intAction {
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
			return
		}
		utils.JsonResponse(ctx, 200, 200200, "处理成功", nil)
		break
	case 0:
		// 取消处理请求
		err1 := services.HandleRequest(intFeedbackID, "null")
		err2 := services.UpdateAdminHaddone(currentUserID, -1)
		// 如果处理过程中出现错误，返回错误信息
		if err1 != nil || err2 != nil {
			utils.LogError(err1)
			utils.LogError(err2)
			utils.JsonResponse(ctx, 200, 200504, "数据库出现在问题，尝试在晚点", nil)
			return
		}
		utils.JsonResponse(ctx, 200, 200200, "取消处理成功", nil)
		break
	default:
		utils.LogError(fmt.Errorf("坏东西来了：未指定的接单操作"))
		utils.JsonResponse(ctx, 200, 200400, "未指定的接单操作", nil)
		break
	}
}

// 评价处理结果（是这个值是写在request里头的，以及要同步更新管理员那里的总分）
// 处理评分请求
func Evaluation(ctx *gin.Context) {
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

	fmt.Println(userType, currentRequest.UndertakerID, currentRequest.IfRubbish, currentRequest.UserID, currentUserID)
	// 如果用户类型为1-且---------请求的执行者id---------不为null-且-请求不是垃圾请求-----------------且----当前用户id----------等于请求的用户id
	if userType == 1 && currentRequest.UndertakerID != "null" && currentRequest.IfRubbish != 0 && currentRequest.UserID == currentUserID {
		// 更新请求的评分和评分内容
		currentRequest.Status = true //改为已处理，最后一步在用户这里
		currentRequest.Grade = input.Grade
		currentRequest.GradeContent = input.GradeContent
		// 更新请求的评分
		if err := services.UpdateRequestEvaluation(&currentRequest); err != nil {
			utils.JsonResponse(ctx, 200, 200403, "评价失败，别灰心，正在让一切重回正轨", nil)
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
		utils.JsonResponse(ctx, 200, 200513, "数据库出现在更新管理员问题", nil)
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
			err = services.RemakeRequest(intID)
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

// 获取指定ID的Request详情
func GetSpecificRequest(ctx *gin.Context) {
	id := ctx.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200516, "错误在反馈ID……", nil)
		return
	}

	request, err1 := services.GetSmallRequestByID(intID)
	replies, err2 := services.GetRepliesByRequestID(intID)
	if err1 != nil && err2 != nil {
		utils.JsonResponse(ctx, 200, 200200, "success", gin.H{
			"request": request,
			"replies": replies,
		})
		return
	} else {
		utils.LogError(err1)
		utils.LogError(err2)
		utils.JsonResponse(ctx, 200, 200517, "好像有点问题，但不大", gin.H{
			"request": request,
			"replies": replies,
		})
		return
	}
}
