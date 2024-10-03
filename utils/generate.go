package utils

import (
	"Back-end/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, userType int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"type":       userType,
		"exp":        time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Config.GetString("jwt.secret")))
	if err != nil {
		LogError(err)
		return ""
	}
	LogError(err)
	return tokenString
}
