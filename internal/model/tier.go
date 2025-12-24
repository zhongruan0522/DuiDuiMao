package model

import "time"

// Tier 额度档位表
type Tier struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`           // 档位名称
	Quota         int       `json:"quota"`          // 额度值
	RequiredLevel int       `json:"required_level"` // 所需信任等级
	DailyLimit    int       `json:"daily_limit"`    // 每人每日限购（0=不限）
	Stock         int       `json:"stock"`          // 当前库存
	IsActive      bool      `json:"is_active"`      // 是否启用
	SortOrder     int       `json:"sort_order"`     // 排序权重
	CreatedAt     time.Time `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`     // 更新时间
}
