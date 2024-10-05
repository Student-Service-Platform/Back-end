package services

import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"
)

// 回复反馈（没有加严格的权限限制，还是要准备为学生回复预留一点的
type requestReply struct {
	request_id string
}

func CheckRequestReplyExistByID(id int64) error {
	result := database.DB.Table("replies").Where("request_id = ?", id).First(&requestReply{})
	if result.Error != nil {
		utils.LogError(result.Error)
	}
	return result.Error
}

func CreateReply(reply *models.Reply) error {
	result := database.DB.Create(reply)
	if result.Error != nil {
		utils.LogError(result.Error)
	}
	return result.Error
}
