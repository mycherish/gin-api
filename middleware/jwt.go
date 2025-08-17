package middleware

import (
	"go-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // 中断请求处理
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) // 响应错误信息
			c.Abort()                                                        // 中断请求处理
			return
		}

		c.Set("user_id", claims.UserID) // 将用户ID保存到上下文中
		c.Next()                        // 继续处理请求
	}
}
