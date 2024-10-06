package utils

import (
	"Back-end/config"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, userType int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"type":       userType,
		"exp":        time.Now().Add(time.Duration(config.Config.GetInt64("cookies.maxAge"))).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Config.GetString("jwt.secret")))
	if err != nil {
		LogError(err)
		return ""
	}
	LogError(err)
	return tokenString
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
