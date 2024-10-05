package services

import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"

	"gorm.io/gorm"
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

// / 不加用户信息，获取所有的
func GetAllRequest(offset int, limit int) ([]formatRequest, error) {
	requests := make([]formatRequest, 0)
	err := database.DB.Offset(offset).Limit(limit).Find(&requests).Error

	if err != nil {
		return nil, err
	}

	//// 处理匿名用户
	for i := range requests {
		if requests[i].IsAnonymous {
			requests[i].Username = "匿名用户"
		}
	}

	return requests, nil
}

// / 加用户信息，获取特定用户的不匿名
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

// 管理员处理帖子的时候同步把UnderTakerID加上
func HandleRequest(requestID int, currentUserID string) error {
	tx := database.DB.Begin()

	if tx.Error != nil {
		utils.LogError(tx.Error)
		return tx.Error
	}
	result := tx.Table("requests").Where("id = ?", requestID).UpdateColumn("undertaker_id", currentUserID)

	if result.Error != nil {
		// 如果更新操作失败，回滚事务
		tx.Rollback()
		utils.LogError(result.Error)
		return result.Error
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// 看一下这个request被处理了没
type HandleInfo struct {
	UndertakerID string `json:"undertaker_id"`
}

func IsHandled(requestID int) (string, error) {
	var handleInfo HandleInfo
	result := database.DB.Table("requests").Where("id = ?", requestID).Find(&handleInfo)
	if result.Error != nil {
		utils.LogError(result.Error)
	}
	return handleInfo.UndertakerID, result.Error
}

// 根据ID获取Request
func GetRequestByID(requestID int) (models.Request, error) {
	var request models.Request
	result := database.DB.Where("id = ?", requestID).Find(&request)
	return request, result.Error
}

// 提交对于回复的评价
func UpdateRequestEvaluation(targetRequst *models.Request) error {
	err := database.DB.Save(targetRequst).Error
	return err
}

// markRequest
func MarkRequest(postID int) error {
	var RubbishRequest models.Request
	result := database.DB.Model(&RubbishRequest).Where("id = ?", postID).Update("if_rubbish", gorm.Expr("if_rubbish + ?", 1))
	return result.Error
}

// remakeRequest
func RemakeRequest(postID int) error {
	var RubbishRequest models.Request
	result := database.DB.Model(&RubbishRequest).Where("id = ?", postID).UpdateColumn("if_rubbish", 0)
	return result.Error
}

// Statue Request
func StatueRequest(postID int) error {
	var RubbishRequest models.Request
	result := database.DB.Model(&RubbishRequest).Where("id = ?", postID).UpdateColumn("if_rubbish", 0)
	return result.Error
}
