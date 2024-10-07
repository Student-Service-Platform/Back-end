package controllers

import (
	"Back-end/utils"

	"github.com/gin-gonic/gin"
)

func BelongsTo(c *gin.Context) {
	userType := c.GetInt("type")

	var t struct {
		UserType int `json:"user_type"`
	}
	c.ShouldBind(&t)

	if t.UserType != userType {
		utils.JsonResponse(c, 200, 200403, "拒绝访问", nil)
		return
	}
	utils.JsonResponse(c, 200, 200200, "允许访问", nil)
}
