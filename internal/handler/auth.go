package handler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/config"
	"github.com/zhongruan0522/DuiDuiMao/internal/service"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	cfg         *config.Config
	userService *service.UserService
	oauthStates map[string]bool // 简易state校验（生产环境应用Redis）
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(cfg *config.Config, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		cfg:         cfg,
		userService: userService,
		oauthStates: make(map[string]bool),
	}
}

// AdminLoginRequest 管理员登录请求（加密后的数据）
type AdminLoginRequest struct {
	Username string `json:"username"` // 双重Base64加密后的用户名
	Password string `json:"password"` // 双重Base64加密后的密码
}

// AdminLogin 管理员账密登录
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "请求参数错误")
		return
	}

	// 解密用户名和密码
	username, err := util.DoubleDecode(req.Username)
	if err != nil {
		util.ErrorResponse(c, 400, "用户名解密失败")
		return
	}

	password, err := util.DoubleDecode(req.Password)
	if err != nil {
		util.ErrorResponse(c, 400, "密码解密失败")
		return
	}

	// 校验管理员账户
	if username != h.cfg.Admin.Username || password != h.cfg.Admin.Password {
		util.ErrorResponse(c, 401, "账号或密码错误")
		return
	}

	// 创建/更新管理员用户记录
	user, err := h.userService.CreateOrUpdateUser(
		0,        // 管理员没有LinuxDoID，使用0
		username, // 使用配置的管理员用户名
		"管理员",    // 昵称
		4,        // 最高信任等级
		true,     // 是管理员
	)
	if err != nil {
		util.ErrorResponse(c, 500, "创建用户记录失败")
		return
	}

	// 生成JWT
	token, err := util.GenerateJWT(user.ID, true, h.cfg.JWT.ExpireHours)
	if err != nil {
		util.ErrorResponse(c, 500, "生成Token失败")
		return
	}

	// 加密响应数据
	encryptedToken := util.DoubleEncode(token)
	encryptedUserData, _ := json.Marshal(user)
	encryptedUser := util.DoubleEncode(string(encryptedUserData))

	util.SuccessResponse(c, gin.H{
		"message": "登录成功",
		"token":   encryptedToken,
		"user":    encryptedUser,
	})
}

// Login LinuxDo OAuth登录跳转
func (h *AuthHandler) Login(c *gin.Context) {
	// 生成随机state
	state := h.generateState()
	h.oauthStates[state] = true

	// 构造授权URL
	authURL := fmt.Sprintf(
		"https://connect.linux.do/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s",
		h.cfg.OAuth.AppClientID,
		url.QueryEscape(h.cfg.OAuth.RedirectURI),
		state,
	)

	util.SuccessResponse(c, gin.H{
		"url": authURL,
	})
}

// Callback LinuxDo OAuth回调
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// 校验state
	if !h.oauthStates[state] {
		util.ErrorResponse(c, 400, "非法的state参数")
		return
	}
	delete(h.oauthStates, state) // 使用后删除

	if code == "" {
		util.ErrorResponse(c, 400, "缺少授权码")
		return
	}

	// 1. 用code换取access_token
	accessToken, err := h.exchangeToken(code)
	if err != nil {
		util.ErrorResponse(c, 500, fmt.Sprintf("获取access_token失败: %v", err))
		return
	}

	// 2. 用access_token获取用户信息
	linuxDoUser, err := h.getLinuxDoUser(accessToken)
	if err != nil {
		util.ErrorResponse(c, 500, fmt.Sprintf("获取用户信息失败: %v", err))
		return
	}

	// 3. 创建/更新本地用户
	user, err := h.userService.CreateOrUpdateUser(
		linuxDoUser.ID,
		linuxDoUser.Username,
		linuxDoUser.Name,
		linuxDoUser.TrustLevel,
		false, // 普通用户不是管理员（除非后续手动设置）
	)
	if err != nil {
		util.ErrorResponse(c, 500, "创建用户记录失败")
		return
	}

	// 4. 生成JWT
	token, err := util.GenerateJWT(user.ID, user.IsAdmin, h.cfg.JWT.ExpireHours)
	if err != nil {
		util.ErrorResponse(c, 500, "生成Token失败")
		return
	}

	// 加密响应数据
	encryptedToken := util.DoubleEncode(token)
	encryptedUserData, _ := json.Marshal(user)
	encryptedUser := util.DoubleEncode(string(encryptedUserData))

	util.SuccessResponse(c, gin.H{
		"message": "登录成功",
		"token":   encryptedToken,
		"user":    encryptedUser,
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	util.SuccessResponse(c, gin.H{
		"message": "登出成功",
	})
}

// ========== 私有辅助方法 ==========

// generateState 生成随机state
func (h *AuthHandler) generateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b) // rand.Read 总是返回 len(b), nil
	return base64.URLEncoding.EncodeToString(b)
}

// exchangeToken 用授权码换取access_token
func (h *AuthHandler) exchangeToken(code string) (string, error) {
	tokenURL := "https://connect.linux.do/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", h.cfg.OAuth.RedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	// 使用Basic Auth
	req.SetBasicAuth(h.cfg.OAuth.AppClientID, h.cfg.OAuth.AppClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// LinuxDoUser LinuxDo用户信息
type LinuxDoUser struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	TrustLevel int    `json:"trust_level"`
	Silenced   bool   `json:"silenced"`
}

// getLinuxDoUser 获取LinuxDo用户信息
func (h *AuthHandler) getLinuxDoUser(accessToken string) (*LinuxDoUser, error) {
	userURL := "https://connect.linux.do/api/user"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var user LinuxDoUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
