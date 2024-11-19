/*
 * @Author: Liu Sainan
 * @Date: 2023-12-07 12:59:22
 */

package server

import (
	"myadmin/internal/config"
	"myadmin/internal/dao"
	"myadmin/internal/db"
	"myadmin/internal/middleware"
	"myadmin/internal/router"
	"myadmin/internal/zlog"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TODO: 初始化 redis (判断是否使用 redis)
// TODO: 初始化 mongodb (判断是否使用 mongodb)

func Server() error {

	config.LoadConfig() // 加载配置文件
	zlog.SetupLogger()  // 初始化 zap 日志
	db.Setup()          // 初始化组件(mysql, redis, cron 等)
	middleware.SetupJWTAuth()

	if err := dao.AutoCreateTable(); err != nil { // 自动创建表结构
		return err
	}

	engine := router.NewRouterEngine()           // 创建路由引擎
	if err := engine.MountRoutes(); err != nil { //挂载路由
		return err
	}

	return runServer(config.GlobalConfig.Server.HttpListenAddress, engine.Engine)
}

func runServer(address string, engine *gin.Engine) error {
	s := endless.NewServer(address, engine)
	s.ReadHeaderTimeout = time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second
	s.ReadTimeout = time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second
	s.WriteTimeout = time.Duration(config.GlobalConfig.Server.ReadTimeout) * time.Second
	s.MaxHeaderBytes = config.GlobalConfig.Server.MaxHeaderBytes << 20

	s.BeforeBegin = func(add string) {
		zap.L().Info("PID 信息", zap.Int("pid", syscall.Getpid()))
	}

	if config.GlobalConfig.Server.UseHttps {
		return s.ListenAndServeTLS(config.GlobalConfig.Server.CertFile, config.GlobalConfig.Server.KeyFile)
	}

	return s.ListenAndServe()
}
