package services

import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"
	"time"
)

// 回复反馈（没有加严格的权限限制，还是要准备为学生回复预留一点的
type requestReply struct {
	request_id string
}

func CreateReply(reply *models.Reply) error {
	result := database.DB.Create(reply)
	if result.Error != nil {
		utils.LogError(result.Error)
	}
	return result.Error
}

// 获取特定反馈的回复
type TinyReply struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	Respondent string    `json:"respondent"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetRepliesByRequestID(request_id int) ([]TinyReply, error) {
	var replies []models.Reply
	result := database.DB.Where("request_id = ? AND father = 0", request_id).Find(&replies)
	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, result.Error
	}

	// 获取所有respondent的ID
	respondentIDs := make([]string, len(replies))
	for i, reply := range replies {
		respondentIDs[i] = reply.Respondent
	}
	// Student列表
	var students []models.Student
	database.DB.Where("id IN ?", respondentIDs).Find(&students)

	// Admin列表
	var admins []models.Admin
	database.DB.Where("id IN ?", respondentIDs).Find(&admins)

	// 创建一个map来存储user_id到username的映射
	usernameMap := make(map[string]string)
	for _, student := range students {
		usernameMap[student.UserID] = student.Username
	}
	for _, admin := range admins {
		usernameMap[admin.UserID] = admin.Username
	}

	// 将respondent替换为username，并转换为TinyReply
	var tinyReplies []TinyReply
	for _, reply := range replies {
		username := reply.Respondent
		if name, exists := usernameMap[reply.Respondent]; exists {
			username = name
		}
		tinyReplies = append(tinyReplies, TinyReply{
			ID:         int(reply.ID),
			Content:    reply.Content,
			Respondent: username,
			CreatedAt:  reply.CreatedAt,
		})
	}

	return tinyReplies, nil
}

// 计算回复数量
func CountRepliesByRequestID(request_id int) (int, error) {
	var count int64
	result := database.DB.Model(&models.Reply{}).Where("request_id = ?", request_id).Count(&count)
	if result.Error != nil {
		utils.LogError(result.Error)
		return 0, result.Error
	}
	return int(count), nil
}
