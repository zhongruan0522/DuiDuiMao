package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// RedeemHandler 兑换处理器
type RedeemHandler struct{}

// NewRedeemHandler 创建兑换处理器
func NewRedeemHandler() *RedeemHandler {
	return &RedeemHandler{}
}

// Redeem 兑换CDK（Mock）
func (h *RedeemHandler) Redeem(c *gin.Context) {
	tierID := c.Param("tier_id")

	util.SuccessResponse(c, gin.H{
		"message":     "兑换成功（Mock数据）",
		"tier_id":     tierID,
		"cdk_code":    "MOCK-CDK-CODE-12345",
		"redeemed_at": time.Now(),
	})
}

// GetHistory 获取兑换记录（Mock）
func (h *RedeemHandler) GetHistory(c *gin.Context) {
	util.SuccessResponse(c, []gin.H{
		{
			"id":          1,
			"tier_name":   "测试档位1",
			"cdk_code":    "MOCK-CDK-001",
			"redeemed_at": time.Now().Add(-24 * time.Hour),
		},
		{
			"id":          2,
			"tier_name":   "测试档位2",
			"cdk_code":    "MOCK-CDK-002",
			"redeemed_at": time.Now().Add(-48 * time.Hour),
		},
	})
}
