package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhongruan/DuiDuiMao/internal/util"
)

// UserHandler 用户处理器
type UserHandler struct{}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetMe 获取当前用户信息（Mock）
func (h *UserHandler) GetMe(c *gin.Context) {
	// 从上下文获取用户ID（由AuthMiddleware设置）
	userID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_admin")

	util.SuccessResponse(c, gin.H{
		"id":          userID,
		"linux_do_id": 12345,
		"username":    "test_user",
		"name":        "测试用户",
		"trust_level": 2,
		"is_admin":    isAdmin,
	})
}
