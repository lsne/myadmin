/*
 * @Author: Liu Sainan
 * @Date: 2023-12-08 23:01:50
 */

package router

import (
	"myadmin/internal/config"
	"myadmin/internal/controller"
	"myadmin/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// func Action(Handler func(*controller.Context)) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := &controller.Context{Context: c}
// 		Handler(ctx)
// 	}
// }

type RouterEngine struct {
	*gin.Engine
}

// 路由引擎配置
func NewRouterEngine() *RouterEngine {

	//上传文件的最大大小,单位M
	maxSize := int64(config.GlobalConfig.Image.MaxSizeMB)

	if config.GlobalConfig.Logger.Level == -1 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.MaxMultipartMemory = maxSize << 20 // 将单位MB转换为Byte

	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.GlobalConfig.Cors.AllowOrigins,
		AllowMethods:     config.GlobalConfig.Cors.AllowMethods,
		AllowHeaders:     config.GlobalConfig.Cors.AllowHeaders,
		ExposeHeaders:    config.GlobalConfig.Cors.ExposeHeaders,
		AllowCredentials: config.GlobalConfig.Cors.AllowCredentials,
	}))

	return &RouterEngine{r}
}

func (r *RouterEngine) MountRoutes() error {

	r.StaticRoutes()

	// root 路由
	root := r.Group(config.GlobalConfig.Server.ApiUrlPrefix)
	{
		user := controller.User{}
		root.POST("/auth/user/register", user.Register)                        //创建用户  TODO: 如果任何人都能注册, 需要在注册时加验证码或者其他手段, 防止攻击脚本疯狂注册导致出现垃圾数据; 如果限制只能管理员添加用户, 则需要将该接口放到 JWTAuth 中间件内部
		root.POST("/auth/user/login", middleware.JWTAuth.Login)                //登录
		root.POST("/auth/user/refresh-token", middleware.JWTAuth.RefreshToken) //刷新 Token
		root.POST("/route/getAsyncRoutes", controller.Route{}.Info)            //刷新 Token
	}

	root.Use(middleware.JWTAuth.LoginRequired)
	// root.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())  // 需要认证和鉴权

	// 加载分组路由
	User(root)

	// 自定义没有路由的处理
	r.NoRoute()

	return nil
}

func (r *RouterEngine) StaticRoutes() {
	// 前端页面
	// TODO: 使用 embed 前端页面打包到 go 二进制程序内部
	r.Static("/favicon.ico", "./dist/favicon.ico")
	r.Static("/assets", "./dist/assets")   // dist里面的静态资源
	r.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	// 为用户头像和文件提供静态地址。
	// TODO: 如果访问文件需要验证, 则需要放到其他分组下
	// TODO: 如果使用 S3 存储文件, 则不需要
	// r.StaticFS("/files/image/avatar/", http.Dir("/data1/file/image/avatar"))
}
