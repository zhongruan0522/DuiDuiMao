package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/service"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// RedeemHandler 兑换处理器
type RedeemHandler struct {
	tierService      *service.TierService
	cdkService       *service.CDKService
	redeemLogService *service.RedeemLogService
}

// NewRedeemHandler 创建兑换处理器
func NewRedeemHandler(tierService *service.TierService, cdkService *service.CDKService, redeemLogService *service.RedeemLogService) *RedeemHandler {
	return &RedeemHandler{
		tierService:      tierService,
		cdkService:       cdkService,
		redeemLogService: redeemLogService,
	}
}

// Redeem 兑换CDK
func (h *RedeemHandler) Redeem(c *gin.Context) {
	tierIDStr := c.Param("tier_id")
	tierID, err := strconv.Atoi(tierIDStr)
	if err != nil {
		util.ErrorResponse(c, 400, "无效的档位ID")
		return
	}

	// 获取登录用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 验证档位是否存在且启用
	tier, err := h.tierService.GetTierByID(tierID)
	if err != nil {
		util.ErrorResponse(c, 404, "档位不存在")
		return
	}
	if !tier.IsActive {
		util.ErrorResponse(c, 400, "该档位未启用")
		return
	}

	// 检查库存
	if tier.Stock <= 0 {
		util.ErrorResponse(c, 400, "该档位已无库存")
		return
	}

	// 获取一个可用的CDK
	cdk, err := h.cdkService.GetAvailableCDKByTierID(tierID)
	if err != nil {
		util.ErrorResponse(c, 500, "获取CDK失败: "+err.Error())
		return
	}

	// 标记CDK为已兑换
	if err := h.cdkService.MarkCDKAsRedeemed(cdk.ID, userID.(int)); err != nil {
		util.ErrorResponse(c, 500, "兑换失败: "+err.Error())
		return
	}

	// 创建兑换记录
	if err := h.redeemLogService.CreateRedeemLog(userID.(int), cdk.ID, tierID); err != nil {
		// 记录失败不影响兑换（CDK已被标记为已兑换）
		// 可以记录日志但不返回错误
	}

	// 解密CDK返回给用户
	cdkCode, err := util.DoubleDecode(cdk.Code)
	if err != nil {
		util.ErrorResponse(c, 500, "CDK解密失败")
		return
	}

	util.SuccessResponse(c, gin.H{
		"message":     "兑换成功",
		"tier_id":     util.DoubleEncode(strconv.Itoa(tierID)),
		"tier_name":   util.DoubleEncode(tier.Name),
		"cdk_code":    cdkCode, // 明文返回给用户
		"redeemed_at": util.DoubleEncode(cdk.RedeemedAt.Format("2006-01-02 15:04:05")),
	})
}

// GetHistory 获取兑换记录
func (h *RedeemHandler) GetHistory(c *gin.Context) {
	// 获取登录用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 获取用户的兑换记录
	redeemLogs, err := h.redeemLogService.GetUserRedeemLogs(userID.(int))
	if err != nil {
		util.ErrorResponse(c, 500, "获取兑换记录失败: "+err.Error())
		return
	}

	// 获取所有档位和CDK用于关联显示
	tiers, _ := h.tierService.GetAllTiers()
	tierMap := make(map[int]string) // tier_id -> tier_name
	for _, tier := range tiers {
		tierMap[tier.ID] = tier.Name
	}

	cdks, _ := h.cdkService.GetCDKs(nil, nil)
	cdkMap := make(map[int]string) // cdk_id -> cdk_code
	for _, cdk := range cdks {
		cdkMap[cdk.ID] = cdk.Code // 已加密
	}

	// 构建返回数据（加密）
	encryptedData := []gin.H{}
	for _, log := range redeemLogs {
		tierName := tierMap[log.TierID]
		cdkCode := cdkMap[log.CDKID]

		// 解密CDK用于显示
		decryptedCode, decErr := util.DoubleDecode(cdkCode)
		if decErr != nil {
			decryptedCode = "***" // 解密失败显示占位符
		}

		encryptedData = append(encryptedData, gin.H{
			"id":          util.DoubleEncode(strconv.Itoa(log.ID)),
			"tier_id":     util.DoubleEncode(strconv.Itoa(log.TierID)),
			"tier_name":   util.DoubleEncode(tierName),
			"cdk_code":    decryptedCode, // 明文显示（已解密）
			"redeemed_at": util.DoubleEncode(log.CreatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	util.SuccessResponse(c, encryptedData)
}
