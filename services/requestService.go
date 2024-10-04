package services

import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"
)

// 发送新Request
func CreateRequest(newRequest models.Request) error {
	result := database.DB.Create(&newRequest)
	return result.Error
}

// 获取Request
type formatRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     int64  `json:"category"`
	Urgency      int    `json:"urgency"`
	IsAnonymous  bool   `json:"is_anonymous"`
	Username     string `json:"username"`
	IfRubbish    int    `json:"if_rubbish"`
	UpdatedAt    string `json:"updated_at"`
	UndertakerID string `json:"undertaker_id"`
	Status       bool   `json:"status"`
}

// 不加用户信息，获取所有的
func GetAllRequest(offset int, limit int) ([]formatRequest, error) {
	requests := make([]formatRequest, 0)
	err := database.DB.Offset(offset).Limit(limit).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	// 处理匿名用户
	for i := range requests {
		if requests[i].IsAnonymous {
			requests[i].Username = "匿名用户"
		}
	}

	return requests, nil
}

// 加用户信息，获取特定用户的不匿名
func GetUserRequest(userid string, offset int, limit int) ([]formatRequest, error) {
	var results []formatRequest

	err := database.DB.Order("-id").Offset(offset).Limit(limit).Where("student_id = ? AND is_anonymous = ?", userid, 0).Find(&results).Error

	if err != nil {
		utils.LogError(err)
		return nil, err
	} else {
		return results, nil
	}
}
