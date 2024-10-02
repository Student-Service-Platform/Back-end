package utils

import (
	"Back-end/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.GetString("jwt.secret")), nil
	})

	if err != nil {
		LogError(err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid")
		fmt.Println("Username:", claims["username"])
		fmt.Println("Type:", int(claims["type"].(float64))) // 注意这里需要将float64转换为int
	} else {
		fmt.Println("Invalid token")
	}
	LogError(err)
	return token, err
}
