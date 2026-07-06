package main

import (
	"fmt"
	"log"

	"github.com/mall-admin/backend/internal/config"
	"github.com/mall-admin/backend/internal/database"
	"github.com/mall-admin/backend/internal/router"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	db, err := database.OpenMySQL(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	engine := router.New(router.Options{
		DB:               db,
		JWTSecret:        cfg.JWT.Secret,
		JWTExpireSeconds: cfg.JWT.ExpireSeconds,
	})
	if err := engine.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
