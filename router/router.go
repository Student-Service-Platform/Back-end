package router

import (
	"Back-end/controllers"
	"Back-end/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	const pre = "/api/auth"

	apiAuth := r.Group(pre) //登录/注册/找回密码
	{
		apiAuth.POST("login", controllers.Login)
		apiAuth.POST("reg", controllers.Register)
	}
	const user = "/api/user" //查看用户信息/修改信息/获取发送的反馈/反馈问题/回复反馈帖子/评价
	apiUser := r.Group(user)
	{
		// 使用JWT中间件保护这些路由
		apiUser.Use(middlewares.TokenAuthMiddleware())
		apiUser.GET("profile", controllers.GetProfile)
		apiUser.PUT("profile", controllers.UpdateProfile)
		// 其他受保护的路由可以在这里添加
	}
}
