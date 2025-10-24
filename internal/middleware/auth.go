package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 webhook 路径
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/webhooks") {
			c.Next()
			return
		}

		// 跳过健康检查路径
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" {
			c.Next()
			return
		}

		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// 检查 Bearer token 格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Invalid token format. Expected 'Bearer <token>'",
				},
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "EMPTY_TOKEN",
					"message": "Token cannot be empty",
				},
			})
			c.Abort()
			return
		}

		// TODO: 验证 token 的有效性
		// 这里应该调用认证服务验证 token
		// 暂时跳过验证，直接设置用户信息
		c.Set("user_id", "system")
		c.Set("user_email", "system@example.com")

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" && tokenParts[1] != "" {
				// TODO: 验证 token 的有效性
				c.Set("user_id", "system")
				c.Set("user_email", "system@example.com")
			}
		}
		c.Next()
	}
}
