package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/config"
	"github.com/zhongruan0522/DuiDuiMao/internal/service"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// AdminHandler 管理端处理器
type AdminHandler struct {
	tierService      *service.TierService
	cdkService       *service.CDKService
	redeemLogService *service.RedeemLogService
}

// NewAdminHandler 创建管理端处理器
func NewAdminHandler(tierService *service.TierService, cdkService *service.CDKService, redeemLogService *service.RedeemLogService) *AdminHandler {
	return &AdminHandler{
		tierService:      tierService,
		cdkService:       cdkService,
		redeemLogService: redeemLogService,
	}
}

// ========== 档位管理 ==========

// TierRequest 档位请求结构
type TierRequest struct {
	Name          string `json:"name" binding:"required"`
	Quota         int    `json:"quota" binding:"required,min=1"`
	RequiredLevel int    `json:"required_level" binding:"min=0,max=4"`
	DailyLimit    int    `json:"daily_limit" binding:"min=0"`
	SortOrder     int    `json:"sort_order"`
	IsActive      bool   `json:"is_active"`
}

// GetTiers 获取档位列表（管理端）
func (h *AdminHandler) GetTiers(c *gin.Context) {
	tiers, err := h.tierService.GetAllTiers()
	if err != nil {
		util.ErrorResponse(c, 500, "获取档位列表失败: "+err.Error())
		return
	}

	// 对数据进行双重Base64加密
	encryptedData := []gin.H{}
	for _, tier := range tiers {
		encryptedData = append(encryptedData, gin.H{
			"id":             util.DoubleEncode(strconv.Itoa(tier.ID)),
			"name":           util.DoubleEncode(tier.Name),
			"quota":          util.DoubleEncode(strconv.Itoa(tier.Quota)),
			"required_level": util.DoubleEncode(strconv.Itoa(tier.RequiredLevel)),
			"daily_limit":    util.DoubleEncode(strconv.Itoa(tier.DailyLimit)),
			"stock":          util.DoubleEncode(strconv.Itoa(tier.Stock)),
			"is_active":      util.DoubleEncode(strconv.FormatBool(tier.IsActive)),
			"sort_order":     util.DoubleEncode(strconv.Itoa(tier.SortOrder)),
			"created_at":     util.DoubleEncode(tier.CreatedAt.Format("2006-01-02 15:04:05")),
			"updated_at":     util.DoubleEncode(tier.UpdatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	util.SuccessResponse(c, encryptedData)
}

// CreateTier 创建档位
func (h *AdminHandler) CreateTier(c *gin.Context) {
	var req TierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	// 创建档位（库存自动计算，无需传入）
	tier, err := h.tierService.CreateTier(
		req.Name,
		req.Quota,
		req.RequiredLevel,
		req.DailyLimit,
		req.SortOrder,
		req.IsActive,
	)
	if err != nil {
		util.ErrorResponse(c, 500, "创建档位失败: "+err.Error())
		return
	}

	// 返回加密后的数据
	util.SuccessResponse(c, gin.H{
		"message": "档位创建成功",
		"id":      util.DoubleEncode(strconv.Itoa(tier.ID)),
	})
}

// UpdateTier 更新档位
func (h *AdminHandler) UpdateTier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.ErrorResponse(c, 400, "无效的档位ID")
		return
	}

	var req TierRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	// 更新档位（库存自动计算，无需传入）
	tier, err := h.tierService.UpdateTier(
		id,
		req.Name,
		req.Quota,
		req.RequiredLevel,
		req.DailyLimit,
		req.SortOrder,
		req.IsActive,
	)
	if err != nil {
		util.ErrorResponse(c, 500, "更新档位失败: "+err.Error())
		return
	}

	util.SuccessResponse(c, gin.H{
		"message": "档位更新成功",
		"id":      util.DoubleEncode(strconv.Itoa(tier.ID)),
	})
}

// DeleteTier 删除档位
func (h *AdminHandler) DeleteTier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.ErrorResponse(c, 400, "无效的档位ID")
		return
	}

	if err := h.tierService.DeleteTier(id); err != nil {
		util.ErrorResponse(c, 500, "删除档位失败: "+err.Error())
		return
	}

	util.SuccessResponse(c, gin.H{
		"message": "档位删除成功",
		"id":      util.DoubleEncode(strconv.Itoa(id)),
	})
}

// ========== CDK管理 ==========

// ImportCDKsRequest 批量导入CDK请求
type ImportCDKsRequest struct {
	TierID string   `json:"tier_id" binding:"required"` // 双重Base64加密的档位ID
	Codes  []string `json:"codes" binding:"required"`   // CDK列表（明文，服务层会加密）
}

// ImportCDKs 批量导入CDK
func (h *AdminHandler) ImportCDKs(c *gin.Context) {
	var req ImportCDKsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	// 解密档位ID
	tierIDStr, err := util.DoubleDecode(req.TierID)
	if err != nil {
		util.ErrorResponse(c, 400, "档位ID解密失败")
		return
	}
	tierID, err := strconv.Atoi(tierIDStr)
	if err != nil {
		util.ErrorResponse(c, 400, "档位ID格式错误")
		return
	}

	// 验证档位是否存在
	tier, err := h.tierService.GetTierByID(tierID)
	if err != nil {
		util.ErrorResponse(c, 404, "档位不存在")
		return
	}

	// 批量导入CDK
	result, err := h.cdkService.BatchImportCDKs(tierID, req.Codes)
	if err != nil {
		util.ErrorResponse(c, 500, "导入CDK失败: "+err.Error())
		return
	}

	// 返回加密后的结果
	util.SuccessResponse(c, gin.H{
		"message":       "CDK导入完成",
		"tier_id":       util.DoubleEncode(strconv.Itoa(tier.ID)),
		"tier_name":     util.DoubleEncode(tier.Name),
		"success_count": util.DoubleEncode(strconv.Itoa(result.SuccessCount)),
		"failed_count":  util.DoubleEncode(strconv.Itoa(result.FailedCount)),
	})
}

// GetCDKs 获取CDK列表
func (h *AdminHandler) GetCDKs(c *gin.Context) {
	// 解析查询参数
	var tierID *int
	var status *int

	if tierIDStr := c.Query("tier_id"); tierIDStr != "" {
		// 解密档位ID
		decryptedTierID, err := util.DoubleDecode(tierIDStr)
		if err != nil {
			util.ErrorResponse(c, 400, "档位ID解密失败")
			return
		}
		tid, err := strconv.Atoi(decryptedTierID)
		if err != nil {
			util.ErrorResponse(c, 400, "档位ID格式错误")
			return
		}
		tierID = &tid
	}

	if statusStr := c.Query("status"); statusStr != "" {
		// 解密状态
		decryptedStatus, err := util.DoubleDecode(statusStr)
		if err != nil {
			util.ErrorResponse(c, 400, "状态解密失败")
			return
		}
		s, err := strconv.Atoi(decryptedStatus)
		if err != nil {
			util.ErrorResponse(c, 400, "状态格式错误")
			return
		}
		status = &s
	}

	// 获取CDK列表
	cdks, err := h.cdkService.GetCDKs(tierID, status)
	if err != nil {
		util.ErrorResponse(c, 500, "获取CDK列表失败: "+err.Error())
		return
	}

	// 获取兑换记录（用于显示兑换用户）
	redeemLogs, _ := h.redeemLogService.GetAllRedeemLogs()
	cdkToLog := make(map[int]int) // cdk_id -> user_id
	for _, log := range redeemLogs {
		cdkToLog[log.CDKID] = log.UserID
	}

	// 对数据进行双重Base64加密
	encryptedData := []gin.H{}
	for _, cdk := range cdks {
		redeemedBy := cdk.RedeemedBy
		if redeemedBy == 0 {
			// 如果CDK表中没有记录，尝试从兑换记录中获取
			if userID, ok := cdkToLog[cdk.ID]; ok {
				redeemedBy = userID
			}
		}

		redeemedAtStr := ""
		if !cdk.RedeemedAt.IsZero() {
			redeemedAtStr = cdk.RedeemedAt.Format("2006-01-02 15:04:05")
		}

		encryptedData = append(encryptedData, gin.H{
			"id":          util.DoubleEncode(strconv.Itoa(cdk.ID)),
			"tier_id":     util.DoubleEncode(strconv.Itoa(cdk.TierID)),
			"code":        cdk.Code, // 已经是加密的
			"status":      util.DoubleEncode(strconv.Itoa(cdk.Status)),
			"redeemed_by": util.DoubleEncode(strconv.Itoa(redeemedBy)),
			"redeemed_at": util.DoubleEncode(redeemedAtStr),
			"created_at":  util.DoubleEncode(cdk.CreatedAt.Format("2006-01-02 15:04:05")),
		})
	}

	util.SuccessResponse(c, encryptedData)
}

// RevokeCDK 作废CDK
func (h *AdminHandler) RevokeCDK(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.ErrorResponse(c, 400, "无效的CDK ID")
		return
	}

	if err := h.cdkService.RevokeCDK(id); err != nil {
		util.ErrorResponse(c, 500, "作废CDK失败: "+err.Error())
		return
	}

	util.SuccessResponse(c, gin.H{
		"message": "CDK作废成功",
		"id":      util.DoubleEncode(strconv.Itoa(id)),
	})
}

// ========== 订单管理（Mock，待后续实现） ==========

// GetOrders 获取订单列表（Mock）
func (h *AdminHandler) GetOrders(c *gin.Context) {
	util.SuccessResponse(c, []gin.H{
		{
			"id":       util.DoubleEncode("1"),
			"user_id":  util.DoubleEncode("1"),
			"username": util.DoubleEncode("test_user"),
			"tier_id":  util.DoubleEncode("1"),
			"quantity": util.DoubleEncode("2"),
			"status":   util.DoubleEncode("1"),
		},
	})
}

// ========== 系统设置 ==========

// GetSettings 获取系统设置
func (h *AdminHandler) GetSettings(c *gin.Context) {
	cfg := config.Get()

	// 返回加密后的数据
	util.SuccessResponse(c, gin.H{
		"global_enabled":       util.DoubleEncode(strconv.FormatBool(cfg.Settings.GlobalEnabled)),
		"announcement":         util.DoubleEncode(cfg.Settings.Announcement),
		"order_expire_minutes": util.DoubleEncode(strconv.Itoa(cfg.Settings.OrderExpireMinutes)),
	})
}

// UpdateSettingsRequest 更新系统设置请求（前端需要加密传输）
type UpdateSettingsRequest struct {
	GlobalEnabled      string `json:"global_enabled"`       // 双重Base64加密的布尔值
	Announcement       string `json:"announcement"`         // 双重Base64加密的公告内容
	OrderExpireMinutes string `json:"order_expire_minutes"` // 双重Base64加密的超时时间
}

// UpdateSettings 更新系统设置
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var encryptedReq UpdateSettingsRequest
	if err := c.ShouldBindJSON(&encryptedReq); err != nil {
		util.ErrorResponse(c, 400, "参数错误: "+err.Error())
		return
	}

	// 解密请求数据
	req := config.UpdateSettingsRequest{}

	if encryptedReq.GlobalEnabled != "" {
		decrypted, err := util.DoubleDecode(encryptedReq.GlobalEnabled)
		if err != nil {
			util.ErrorResponse(c, 400, "global_enabled 解密失败")
			return
		}
		enabled := decrypted == "true"
		req.GlobalEnabled = &enabled
	}

	if encryptedReq.Announcement != "" {
		decrypted, err := util.DoubleDecode(encryptedReq.Announcement)
		if err != nil {
			util.ErrorResponse(c, 400, "announcement 解密失败")
			return
		}
		req.Announcement = &decrypted
	}

	if encryptedReq.OrderExpireMinutes != "" {
		decrypted, err := util.DoubleDecode(encryptedReq.OrderExpireMinutes)
		if err != nil {
			util.ErrorResponse(c, 400, "order_expire_minutes 解密失败")
			return
		}
		minutes, err := strconv.Atoi(decrypted)
		if err != nil {
			util.ErrorResponse(c, 400, "order_expire_minutes 格式错误")
			return
		}
		if minutes < 1 {
			util.ErrorResponse(c, 400, "订单超时时间必须大于0")
			return
		}
		req.OrderExpireMinutes = &minutes
	}

	// 更新配置文件并热重载
	if err := config.UpdateSettings(req); err != nil {
		util.ErrorResponse(c, 500, "更新设置失败: "+err.Error())
		return
	}

	util.SuccessResponse(c, gin.H{
		"message": "设置更新成功",
	})
}
