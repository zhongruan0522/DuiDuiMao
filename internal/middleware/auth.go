package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header或Cookie获取token
		token := c.GetHeader("Authorization")
		if token != "" {
			// 去除Bearer前缀
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			// 尝试从Cookie获取
			token, _ = c.Cookie("token")
		}

		if token == "" {
			util.ErrorResponse(c, 401, "未登录")
			c.Abort()
			return
		}

		// 解析JWT
		claims, err := util.ParseJWT(token)
		if err != nil {
			util.ErrorResponse(c, 401, "Token无效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}
