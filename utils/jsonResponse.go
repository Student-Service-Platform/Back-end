package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 基本的响应模板
func JsonResponse(ctx *gin.Context, httpStatusCode int, code int, msg string, data interface{}) {
	ctx.JSON(httpStatusCode, gin.H{
		"code": code, //六位的代码，应该是200+XXX的格式，如果没有出错就200，出错一般200+XXX
		"data": data, //返回的数据
		"msg":  msg,  //返回的消息，默认是success
	})
}

// 成功操作
func JsonSuccessResponse(ctx *gin.Context, data interface{}) {
	JsonResponse(ctx, http.StatusOK, 200200, "success", data)
}

// 未授权的访问
func JsonUnAuthorizedResponse(ctx *gin.Context, msg string) {
	JsonResponse(ctx, http.StatusUnauthorized, 200401, msg, nil)
	ctx.Abort() // 终止后续处理
}
