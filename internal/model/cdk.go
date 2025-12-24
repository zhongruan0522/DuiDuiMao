package model

import "time"

// CDK CDK表
type CDK struct {
	ID         int       `json:"id"`
	TierID     int       `json:"tier_id"`     // 所属档位ID
	Code       string    `json:"code"`        // CDK内容（加密存储）
	Status     int       `json:"status"`      // 0:未兑换 1:已锁定 2:已兑换 3:已作废
	OrderID    int       `json:"order_id"`    // 关联订单ID（0表示未关联）
	RedeemedBy int       `json:"redeemed_by"` // 兑换用户ID（0表示未兑换）
	RedeemedAt time.Time `json:"redeemed_at"` // 兑换时间
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}
