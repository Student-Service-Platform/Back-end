package controllers

import (
	"Back-end/services"
	"Back-end/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//管理员注册

type AdminRegister struct {
	AdminID  string `json:"admin_id"   binding:"required"`
	Username string `json:"username"   binding:"required"`
	Password string `json:"password"  binding:"required"`
	MailAuth bool   `json:"mail_auth"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
}

func Admin_Register(ctx *gin.Context) {
	if ctx.GetInt("type") != 3 { // 筛选超管请求
		utils.JsonResponse(ctx, 401, 401, "权限不足", nil)
		return
	}

	var data AdminRegister
	// 解析参数并且绑定到 struct 中
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		utils.LogError(err)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}

	// 检查管理员id是否合法（A+字符串）

	if data.AdminID[0] != 'a' {
		utils.JsonResponse(ctx, 200, 200501, "你这id有问题啊", nil)
		utils.LogError(err)
	} else { //确定工号没问题之后检验密码
		if len(data.Password) < 8 || len(data.Password) > 16 {
			utils.JsonResponse(ctx, 200, 200502, "你这密码有问题啊（密码过长或过短）", nil)
		} else {
			//检验是否存在
			err = services.CheckUserExistByUserID(data.AdminID, "admins")
			if err != nil && err != gorm.ErrRecordNotFound { //如果发生了错误
				utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
			} else if err == nil { //如果找到了
				utils.JsonResponse(ctx, 200, 200504, "你这id已经被注册了", nil)
			} else {
				//补齐struct然后塞进去，这里只负责传入相关信息
				err = services.AddAdmin(data.AdminID, data.Username, data.Password)
				if err != nil {
					utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
					utils.LogError(err)
				} else {
					utils.JsonResponse(ctx, 200, 200200, "注册成功", nil)
				}
			}
		}
	}
}

func Del(ctx *gin.Context) {
	if ctx.GetInt("type") != 3 { // 筛选超管请求
		utils.JsonResponse(ctx, 401, 401, "权限不足", nil)
		return
	}

	DelUser := ctx.Query("UserID")
	if DelUser == "" {
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}
	var table string
	if DelUser[0] == 'a' {
		table = "admins"
	} else {
		table = "students"
	}
	err := services.CheckUserExistByUserID(DelUser, table)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
	} else if err == nil { //如果找到了
		err = services.DelUser(DelUser, table)
		if err != nil {
			utils.JsonResponse(ctx, 200, 200503, "你今天有点问题啊(bushi)，咱遇到了点问题，晚点再试吧", nil)
			utils.LogError(err)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "删除成功", nil)
		}
	}
}

func GetRubbish(ctx *gin.Context) {
	if ctx.GetInt("type") != 3 { // 筛选超管请求
		utils.JsonResponse(ctx, 401, 401, "权限不足", nil)
		return
	}

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
	offset := (page - 1) * perPage // 计算偏移量             //可以看得到所有的Request，包括匿名的

	// 获取所有的垃圾信息
	requests, err := services.GetAllRubbish(offset, perPage)
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

// 审批垃圾反馈

func UpdateRubbish(ctx *gin.Context) {
	if ctx.GetInt("type") != 3 { // 筛选超管请求
		utils.JsonResponse(ctx, 401, 401, "权限不足", nil)
		return
	}

	strID := ctx.Query("id")
	strAction := ctx.Query("action")

	intID, err1 := strconv.Atoi(strID)
	intAction, err2 := strconv.Atoi(strAction)

	if err1 != nil || err2 != nil {
		utils.LogError(err1)
		utils.LogError(err2)
		utils.JsonResponse(ctx, 200, 200503, "参数错误", nil)
		return
	}
	switch intAction {
	case 0:
		err := services.RemakeRequest(intID)
		if err != nil {
			utils.JsonResponse(ctx, 200, 200503, "今天看来不太好！等等再试吧", nil)
			utils.LogError(err)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "操作成功", nil)
		}
		break
	case 1:
		err := services.RubbishRequest(intID)
		if err != nil {
			utils.JsonResponse(ctx, 200, 200503, "今天看来不太好！等等再试吧", nil)
			utils.LogError(err)
		} else {
			utils.JsonResponse(ctx, 200, 200200, "操作成功", nil)
		}
		break
	default:
		utils.JsonResponse(ctx, 200, 200503, "你在干一些非法的操作", nil)
		break
	}
}

func GetUandA(ctx *gin.Context) {
	table := ctx.Query("type")
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
	offset := (page - 1) * perPage // 计算偏移量             //可以看得到所有的Request，包括匿名的

	// 获取所有的垃圾信息
	requests, err := services.GetAll(table, offset, perPage)
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
