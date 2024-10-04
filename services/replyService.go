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

var NewReply models.Reply

func CreateRequestReply(content string, currentUserID string, request_id int64) error {
	NewReply = models.Reply{
		RequestID:  request_id,
		Father:     0,
		Content:    content,
		Respondent: currentUserID,
	}

	result := database.DB.Table("replies").Create(&NewReply)
	return result.Error
}
