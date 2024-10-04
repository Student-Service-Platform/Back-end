package services

import (
	"Back-end/database"
	"Back-end/models"
)

// 发送新Request
func CreateRequest(newRequest models.Request) error {
	result := database.DB.Create(&newRequest)
	return result.Error
}
