/*
 * @Author: Liu Sainan
 * @Date: 2023-12-09 02:00:59
 */

package router

import (
	"myadmin/internal/controller"
	"myadmin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func User(root *gin.RouterGroup) {
	user := controller.User{}
	userRouter := root.Group("/auth", func(ctx *gin.Context) {})
	{
		userRouter.POST("/user/logout", middleware.JWTAuth.Logout)                  //退出登录
		userRouter.POST("/user/info", user.Info)                                    //获取个人用户信息
		userRouter.POST("/user/users", middleware.JWTAuth.RootRequired, user.Users) //获取用户列表
		userRouter.POST("/user/update", user.Update)
		userRouter.POST("/user/delete", middleware.JWTAuth.RootRequired, user.Delete)
		userRouter.POST("/user/password/update", user.UpdatePassword)
		userRouter.POST("/user/password/reset", middleware.JWTAuth.RootRequired, user.ResetPassword)
		// 添加更新头像的接口
	}
}
