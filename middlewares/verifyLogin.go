package middlewares

import (
	"Back-end/config"
	"Back-end/services"
	"Back-end/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// VerifyJWT 函数用于验证 JWT 的合法性
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	// 解析 JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回用于验证签名的密钥
		return []byte(config.Config.GetString("jwt.secret")), nil
	})

	// 检查解析过程中是否出现错误
	if err != nil {
		utils.LogError(err)
		return nil, err
	}

	// 检查 JWT 是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid")
		fmt.Println("UserID:", claims["user_id"])
		fmt.Println("Type:", int(claims["type"].(float64))) // 注意这里需要将float64转换为int
	} else {
		fmt.Println("Invalid token")
	}
	utils.LogError(err)
	return token, err
}

// TokenAuthMiddleware 是一个中间件函数，用于验证 JWT 的合法性
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Received Token:", tokenString) // 打印接收到的token
		if tokenString == "" {
			// 如果 Authorization 字段为空，返回 401 Unauthorized 错误
			c.JSON(http.StatusUnauthorized, gin.H{"code": 200401, "data": nil, "msg": "未登录"})
			c.Abort()
			return
		}

		// 验证 JWT 的合法性
		token, err := VerifyJWT(tokenString)
		if err != nil || !token.Valid {
			// 如果 JWT 验证失败，返回 401 Unauthorized 错误
			utils.LogError(err)
			c.JSON(http.StatusUnauthorized, gin.H{"code": 200401, "data": nil, "msg": "登录信息校验失败"})
			c.Abort()
			return
		}

		// 检查 JWT 的 claims 是否有效
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// 如果 JWT 的 claims 无效，返回 401 Unauthorized 错误
			c.JSON(http.StatusUnauthorized, gin.H{"code": 200401, "data": nil, "msg": "登录信息（claims）无效"})
			c.Abort()
			return
		}

		// 检查 user_id 是否存在
		userID := claims["user_id"].(string)
		userType := int(claims["type"].(float64))
		var table string
		switch userType {
		case 1:
			table = "students"
			break
		case 2, 3:
			table = "admins"
			break
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"code": 200401, "data": nil, "msg": "登录信息（用户类型）无效"})
			c.Abort()
			return
		}

		if services.CheckUserExistByUserID(userID, table) != nil {
			// 如果 user_id 不存在，返回 401 Unauthorized 错误
			c.JSON(http.StatusUnauthorized, gin.H{"code": 200401, "data": nil, "msg": "登录信息（用户ID）无效"})
			utils.LogError(err)
			c.Abort()
			return
		}

		// 将解析得到的 userID 和 type 参数传递给下一步命令
		c.Set("userID", claims["user_id"].(string))
		c.Set("type", int(claims["type"].(float64)))

		// 继续处理请求
		c.Next()
	}
}
