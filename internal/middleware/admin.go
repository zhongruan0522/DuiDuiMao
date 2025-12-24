package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhongruan/DuiDuiMao/internal/util"
)

// AdminMiddleware 管理员权限中间件（需要先经过AuthMiddleware）
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			util.ErrorResponse(c, 403, "无管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
