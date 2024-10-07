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

	const user = "/api/user" //查看用户信息/修改信息/反馈问题/回复反馈帖子/评价
	apiUser := r.Group(user)
	{
		// 使用JWT中间件保护这些路由
		apiUser.Use(middlewares.TokenAuthMiddleware())
		apiUser.GET("profile", controllers.GetProfile)
		apiUser.PUT("profile", controllers.UpdateProfile)
		apiUser.POST("feedback", controllers.CreateRequest)
		// 其他受保护的路由可以在这里添加
	}
	r.GET("/api/user/feedback", controllers.GetAllRequests)
	r.GET("/api/feedback/:id", controllers.GetSpecificRequest)
	r.GET("/api/feedback/select", controllers.GetSelectedFeedback)

	const feedback = "/api/feedback"
	apiFeedback := r.Group(feedback)
	{
		apiFeedback.Use(middlewares.TokenAuthMiddleware(), middlewares.ValidPath()) //
		apiFeedback.POST(":id/reply", controllers.ReplyRequest)
		apiFeedback.PUT(":id/admin", controllers.HandleRequest)
		apiFeedback.PUT(":id/evaluation", controllers.Evaluation)
		apiFeedback.PUT(":id/mark", controllers.MarkRequest)
	}

	const superadmin = "/api/superadmin"
	apiSuperAdmin := r.Group(superadmin)
	{
		// apiSuperAdmin.POST("")
		apiSuperAdmin.Use(middlewares.TokenAuthMiddleware(), middlewares.ValidPath()) //
		apiSuperAdmin.GET("U&A", controllers.GetUandA)
		apiSuperAdmin.POST("workbench", controllers.Admin_Register)
		apiSuperAdmin.DELETE("workbench", controllers.Del)
		apiSuperAdmin.GET("workbench", controllers.GetRubbish)
		apiSuperAdmin.PUT("workbench", controllers.UpdateRubbish)

	}
	apiGuard := r.Group("/api/guard")
	{
		apiGuard.Use(middlewares.TokenAuthMiddleware(), middlewares.ValidPath())
		apiGuard.POST("/", controllers.BelongsTo)
	}

}
