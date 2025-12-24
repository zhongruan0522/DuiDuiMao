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
	tierService *service.TierService
}

// NewAdminHandler 创建管理端处理器
func NewAdminHandler(tierService *service.TierService) *AdminHandler {
	return &AdminHandler{
		tierService: tierService,
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
	if err := c.ShouldBindJSON(&req); err != nil {
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

// ========== CDK管理（Mock，待后续实现） ==========

// ImportCDKs 批量导入CDK（Mock）
func (h *AdminHandler) ImportCDKs(c *gin.Context) {
	util.SuccessResponse(c, gin.H{
		"message": "CDK导入成功（Mock）",
		"count":   util.DoubleEncode("100"),
	})
}

// GetCDKs 获取CDK列表（Mock）
func (h *AdminHandler) GetCDKs(c *gin.Context) {
	util.SuccessResponse(c, []gin.H{
		{
			"id":      util.DoubleEncode("1"),
			"tier_id": util.DoubleEncode("1"),
			"code":    util.DoubleEncode("MOCK-CODE-001"),
			"status":  util.DoubleEncode("0"),
		},
	})
}

// RevokeCDK 作废CDK（Mock）
func (h *AdminHandler) RevokeCDK(c *gin.Context) {
	id := c.Param("id")
	util.SuccessResponse(c, gin.H{
		"message": "CDK作废成功（Mock）",
		"id":      util.DoubleEncode(id),
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
