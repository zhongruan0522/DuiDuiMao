package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 运行模式常量
const (
	ModeDev    = "dev"    // 开发模式
	ModeServer = "server" // 服务器模式
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	OAuth    OAuthConfig
	Pay      PayConfig
	Admin    AdminConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Settings SettingsConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
	Mode string // dev 或 server
}

// OAuthConfig OAuth认证配置
type OAuthConfig struct {
	AppClientID     string
	AppClientSecret string
	RedirectURI     string
}

// PayConfig 支付配置
type PayConfig struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	NotifyURL    string
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username string
	Password string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	URL string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// SettingsConfig 系统设置配置
type SettingsConfig struct {
	GlobalEnabled      bool   // 全局开关（是否暂停购买功能）
	Announcement       string // 公告内容
	OrderExpireMinutes int    // 订单超时时间（分钟）
}

var (
	globalConfig *Config
	configMutex  sync.RWMutex
)

// Load 从config.yaml加载配置（使用简单的键值对格式）
func Load() (*Config, error) {
	cfg, err := loadFromFile()
	if err != nil {
		return nil, err
	}

	// 保存全局配置
	configMutex.Lock()
	globalConfig = cfg
	configMutex.Unlock()

	return cfg, nil
}

// loadFromFile 从文件加载配置
func loadFromFile() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	cfg := &Config{}
	lines := strings.Split(string(data), ",")

	configMap := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		configMap[key] = value
	}

	// 解析配置
	cfg.Server.Port, _ = strconv.Atoi(getConfigValue(configMap, "port", "3001"))
	cfg.Server.Mode = getConfigValue(configMap, "mode", ModeDev)

	cfg.OAuth.AppClientID = getConfigValue(configMap, "app_client_id", "")
	cfg.OAuth.AppClientSecret = getConfigValue(configMap, "app_client_secret", "")
	cfg.OAuth.RedirectURI = getConfigValue(configMap, "app_redirect_uri", "http://localhost:3001/api/auth/callback")

	cfg.Pay.ClientID = getConfigValue(configMap, "pay_client_id", "")
	cfg.Pay.ClientSecret = getConfigValue(configMap, "pay_client_secret", "")
	cfg.Pay.CallbackURL = getConfigValue(configMap, "pay_callback_url", "http://localhost:3001/api/pay/callback")
	cfg.Pay.NotifyURL = getConfigValue(configMap, "pay_notify_url", "http://localhost:3001/api/pay/notify")

	cfg.Admin.Username = getConfigValue(configMap, "adminname", "root")
	cfg.Admin.Password = getConfigValue(configMap, "adminpassword", "rootpassword")

	cfg.Database.URL = getConfigValue(configMap, "dburl", "")

	cfg.JWT.Secret = getConfigValue(configMap, "jwt_secret", "default_jwt_secret")
	cfg.JWT.ExpireHours, _ = strconv.Atoi(getConfigValue(configMap, "jwt_expire_hours", "168"))

	// 解析系统设置
	cfg.Settings.GlobalEnabled = getConfigValue(configMap, "global_enabled", "true") == "true"
	cfg.Settings.Announcement = getConfigValue(configMap, "announcement", "欢迎使用兑兑猫 CDK 兑换平台！")
	cfg.Settings.OrderExpireMinutes, _ = strconv.Atoi(getConfigValue(configMap, "order_expire_minutes", "15"))

	return cfg, nil
}

func getConfigValue(m map[string]string, key, defaultValue string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return defaultValue
}

// Get 获取当前全局配置（线程安全）
func Get() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// Reload 重新加载配置文件（热重载）
func Reload() error {
	cfg, err := loadFromFile()
	if err != nil {
		return err
	}

	configMutex.Lock()
	globalConfig = cfg
	configMutex.Unlock()

	return nil
}

// UpdateSettingsRequest 更新系统设置请求
type UpdateSettingsRequest struct {
	GlobalEnabled      *bool   `json:"global_enabled"`
	Announcement       *string `json:"announcement"`
	OrderExpireMinutes *int    `json:"order_expire_minutes"`
}

// UpdateSettings 更新系统设置并写入配置文件（热更新）
func UpdateSettings(req UpdateSettingsRequest) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// 读取当前配置文件内容
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	configMap := make(map[string]string)
	lines := strings.Split(string(data), ",")

	// 解析现有配置
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		configMap[key] = value
	}

	// 更新设置
	if req.GlobalEnabled != nil {
		configMap["global_enabled"] = strconv.FormatBool(*req.GlobalEnabled)
		globalConfig.Settings.GlobalEnabled = *req.GlobalEnabled
	}
	if req.Announcement != nil {
		configMap["announcement"] = *req.Announcement
		globalConfig.Settings.Announcement = *req.Announcement
	}
	if req.OrderExpireMinutes != nil {
		configMap["order_expire_minutes"] = strconv.Itoa(*req.OrderExpireMinutes)
		globalConfig.Settings.OrderExpireMinutes = *req.OrderExpireMinutes
	}

	// 写回文件
	var newLines []string
	for key, value := range configMap {
		newLines = append(newLines, fmt.Sprintf("\"%s\"=\"%s\"", key, value))
	}
	newContent := strings.Join(newLines, ",\n")

	if err := os.WriteFile("config.yaml", []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
