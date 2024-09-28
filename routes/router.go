package routes

import "github.com/gin-gonic/gin"

func CollectRoutes(r *gin.Engine) *gin.Engine {
	v1 := r.Group("/api/")
	{

	}
	return v1
}
