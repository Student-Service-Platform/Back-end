package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateID 中间件用于验证路径参数 id
func ValidPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}
