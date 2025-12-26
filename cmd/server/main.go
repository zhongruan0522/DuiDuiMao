package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zhongruan0522/DuiDuiMao/internal/config"
	"github.com/zhongruan0522/DuiDuiMao/internal/handler"
	"github.com/zhongruan0522/DuiDuiMao/internal/middleware"
	"github.com/zhongruan0522/DuiDuiMao/internal/service"
	"github.com/zhongruan0522/DuiDuiMao/internal/util"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化JWT
	util.InitJWT(cfg.JWT.Secret)

	// 创建Gin引擎
	r := gin.Default()

	// 创建服务层
	userService := service.NewUserService(cfg.Server.Mode)
	tierService := service.NewTierService(cfg.Server.Mode)
	cdkService := service.NewCDKService(cfg.Server.Mode)
	redeemLogService := service.NewRedeemLogService(cfg.Server.Mode)

	// 创建处理器
	authHandler := handler.NewAuthHandler(cfg, userService)
	userHandler := handler.NewUserHandler()
	redeemHandler := handler.NewRedeemHandler(tierService, cdkService, redeemLogService)
	adminHandler := handler.NewAdminHandler(tierService, cdkService, redeemLogService)

	// ========== 用户端接口 ==========
	api := r.Group("/api")
	{
		// 认证接口（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/admin/login", authHandler.AdminLogin) // 管理员账密登录
			auth.GET("/login", authHandler.Login)             // LinuxDo OAuth登录跳转
			auth.GET("/callback", authHandler.Callback)       // OAuth回调
			auth.POST("/logout", authHandler.Logout)          // 登出
		}

		// 用户接口（需要登录）
		user := api.Group("/user", middleware.AuthMiddleware())
		{
			user.GET("/me", userHandler.GetMe)
		}

		// 档位接口（无需登录）
		api.GET("/tiers", adminHandler.GetTiers) // 暂时用管理端的接口，后续可以创建用户端专用接口

		// 兑换接口（需要登录）
		redeem := api.Group("/redeem", middleware.AuthMiddleware())
		{
			redeem.POST("/:tier_id", redeemHandler.Redeem)
			redeem.GET("/history", redeemHandler.GetHistory)
		}

		// ========== 管理端接口 ==========
		admin := api.Group("/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// 档位管理
			admin.GET("/tiers", adminHandler.GetTiers)
			admin.POST("/tiers", adminHandler.CreateTier)
			admin.PUT("/tiers/:id", adminHandler.UpdateTier)
			admin.DELETE("/tiers/:id", adminHandler.DeleteTier)

			// CDK管理
			admin.POST("/cdks/import", adminHandler.ImportCDKs)
			admin.GET("/cdks", adminHandler.GetCDKs)
			admin.PUT("/cdks/:id/revoke", adminHandler.RevokeCDK)

			// 订单管理
			admin.GET("/orders", adminHandler.GetOrders)

			// 系统设置
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)
		}
	}

	// 静态文件服务（前端）
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 DuiDuiMao 服务启动成功！")
	log.Printf("📍 监听地址: http://localhost%s", addr)
	log.Printf("🎯 运行模式: %s", cfg.Server.Mode)

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
