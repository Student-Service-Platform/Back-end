package services

import (
	"Back-end/database"
	"Back-end/models"
	"Back-end/utils"
	"time"

	"gorm.io/gorm"
)

// 发送新Request
func CreateRequest(newRequest models.Request) error {
	result := database.DB.Create(&newRequest)
	return result.Error
}

// 获取Request
// 定义一个结构体RequestInfo，用于存储请求信息
type RequestInfo struct {
	Username    string    `json:"username"`    // 请求者用户名
	CreatedAt   time.Time `json:"created_at"`  // 请求创建时间
	Title       string    `json:"title"`       // 请求标题
	Description string    `json:"description"` // 请求描述
	Category    int64     `json:"category"`    // 请求类别
	Urgency     int64     `json:"urgency"`     // 请求紧急程度
	IfRubbish   int64     `json:"if_rubbish"`  // 是否为垃圾请求
	Undertaker  string    `json:"undertaker"`  // 负责人用户名
	Status      bool      `json:"status"`      // 请求状态
}

// 获取请求信息
// 根据偏移量和限制获取请求信息
func GetAllRequests(offset, limit int) ([]RequestInfo, error) {
	var requests []models.Request
	if err := database.DB.Offset(offset).Limit(limit).Find(&requests).Error; err != nil {
		return nil, err
	}

	var requestInfos []RequestInfo
	for _, req := range requests {
		var student models.Student
		if err := database.DB.Where("user_id = ?", req.UserID).First(&student).Error; err != nil {
			return nil, err
		}

		var admin models.Admin
		if req.UndertakerID != "" {
			if err := database.DB.Where("user_id = ?", req.UndertakerID).First(&admin).Error; err != nil {
				return nil, err
			}
		}

		username := student.Username
		if req.IsAnonymous {
			username = "匿名用户"
		}

		undertaker := ""
		if req.UndertakerID != "" {
			undertaker = admin.Username
		}

		requestInfo := RequestInfo{
			Username:    username,
			CreatedAt:   req.CreatedAt,
			Title:       req.Title,
			Description: req.Description,
			Category:    req.Category,
			Urgency:     req.Urgency,
			IfRubbish:   req.IfRubbish,
			Undertaker:  undertaker,
			Status:      req.Status,
		}
		requestInfos = append(requestInfos, requestInfo)
	}

	return requestInfos, nil
}

func GetRequestsByUserID(userID string, offset, limit int) ([]RequestInfo, error) {
	var requests []models.Request
	if err := database.DB.Where("user_id = ? AND is_anonymous = ?", userID, false).Offset(offset).Limit(limit).Find(&requests).Error; err != nil {
		return nil, err
	}

	var requestInfos []RequestInfo
	for _, req := range requests {
		var student models.Student
		if err := database.DB.Where("user_id = ?", userID).First(&student).Error; err != nil {
			return nil, err
		}

		var admin models.Admin
		if req.UndertakerID != "0" {
			if err := database.DB.Where("user_id = ?", req.UndertakerID).First(&admin).Error; err != nil {
				return nil, err
			}
		}

		undertaker := ""
		if req.UndertakerID != "0" {
			undertaker = admin.Username
		}

		requestInfo := RequestInfo{
			Username:    student.Username,
			CreatedAt:   req.CreatedAt,
			Title:       req.Title,
			Description: req.Description,
			Category:    req.Category,
			Urgency:     req.Urgency,
			IfRubbish:   req.IfRubbish,
			Undertaker:  undertaker,
			Status:      req.Status,
		}
		requestInfos = append(requestInfos, requestInfo)
	}

	return requestInfos, nil
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
