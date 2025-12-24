package model

import "time"

// RedeemLog 兑换记录表
type RedeemLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CDKID     int       `json:"cdk_id"`
	TierID    int       `json:"tier_id"` // 所属档位ID
	CreatedAt time.Time `json:"created_at"`
}
