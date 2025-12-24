package model

import "time"

// User 用户表
type User struct {
	ID         int       `json:"id"`
	LinuxDoID  int       `json:"linux_do_id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	TrustLevel int       `json:"trust_level"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
