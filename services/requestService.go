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
	Category    int       `json:"category"`    // 请求类别
	Urgency     int       `json:"urgency"`     // 请求紧急程度
	IfRubbish   int       `json:"if_rubbish"`  // 是否为垃圾请求
	Undertaker  string    `json:"undertaker"`  // 负责人用户名
	Status      bool      `json:"status"`      // 请求状态
}

// GetAllRequests retrieves all requests from the database, with an optional offset and limit.
func GetAllRequests(offset, limit int) ([]RequestInfo, error) {
	// Define a variable to store the requests
	var requests []models.Request
	// Retrieve requests from the database with an offset and limit
	if err := database.DB.Offset(offset).Limit(limit).
		// Preload the student information
		Preload("Student", "user_id = ?", "user_id").
		// Preload the admin information
		Preload("Admin", "user_id = ?", "undertaker_id").
		// Find the requests
		Find(&requests).Error; err != nil {
		// Return an error if the request fails
		return nil, err
	}

	// Define a variable to store the request information
	var requestInfos []RequestInfo
	// Iterate through the requests
	for _, req := range requests {
		// Set the username to the student's username
		username := req.Student.Username
		// If the request is anonymous, set the username to "匿名用户"
		if req.IsAnonymous {
			username = "匿名用户"
		}

		// Set the undertaker to an empty string
		undertaker := ""
		// If the undertaker ID is not "null", set the undertaker to the admin's username
		if "null" != req.UndertakerID {
			undertaker = req.Admin.Username
		}

		// Create a new request information object
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
		// Append the request information to the request information array
		requestInfos = append(requestInfos, requestInfo)
	}

	// Return the request information array
	return requestInfos, nil
}

func GetRequestsByUserID(targetUserID string, offset, limit int) ([]RequestInfo, error) {
	var requests []models.Request
	if err := database.DB.Where("user_id = ? AND is_anonymous = ?", targetUserID, false).Offset(offset).Limit(limit).Find(&requests).Error; err != nil {
		utils.LogError(err)
		return nil, err
	}

	var targetUser models.Student
	if err := database.DB.Where("user_id = ?", targetUserID).First(&targetUser).Error; err != nil {
		utils.LogError(err)
		return nil, err
	}

	var requestInfos []RequestInfo
	for _, req := range requests {

		var admin models.Admin
		if req.UndertakerID != "null" {
			if err := database.DB.Where("user_id = ?", req.UndertakerID).First(&admin).Error; err != nil {
				return nil, err
			}
		}

		undertaker := ""
		if req.UndertakerID != "null" {
			undertaker = admin.Username
		}

		requestInfo := RequestInfo{
			Username:    targetUser.Username,
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
