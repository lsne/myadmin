/*
 * @Author: Liu Sainan
 * @Date: 2024-01-06 23:28:25
 */

package controller

import (
	"myadmin/internal/config"
	"myadmin/internal/dto"
	"myadmin/internal/middleware"
	"myadmin/internal/model"
	"myadmin/internal/service/userservice"
	"myadmin/internal/utils/ginutils"

	"github.com/gin-gonic/gin"
)

type User struct {
}

func (u User) Register(c *gin.Context) {

	var req dto.UserRegisterReq

	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	if err := req.AvoidXSS(); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	if err := req.Validator(); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	userService := userservice.NewUserService()

	if err := userService.Register(req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	ginutils.RespOK(c, "注册成功")
}

func (u User) Info(c *gin.Context) {

	v, ok := c.Get(config.GinCtxUserKey)
	if !ok {
		ginutils.RespError(c, "未找到用户信息")
	}

	switch user := v.(type) {
	case nil:
		ginutils.RespError(c, "未找到用户信息")
	case model.User:
		ginutils.RespData(c, user)
	default:
		ginutils.RespError(c, "解析用户信息失败")
	}
}

func (u User) Users(c *gin.Context) {
	var req dto.UserListReq
	var resp dto.UserListResp
	var err error

	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	userService := userservice.NewUserService()
	if resp, err = userService.Users(req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}
	ginutils.RespData(c, resp)
}

func (u User) Update(c *gin.Context) {
	var req dto.UserUpdateReq
	var err error

	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	if err := middleware.JWTAuth.RootOrOwnerRequired(c, req.ID); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	if err := req.AvoidXSS(); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	if err := req.Validator(); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	userService := userservice.NewUserService()
	if err = userService.Update(req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	ginutils.RespOK(c, "用户更新成功")
}

func (u User) UpdatePassword(c *gin.Context) {
	var req dto.UserUpdatePasswordReq
	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	var user model.User

	v, ok := c.Get(config.GinCtxUserKey)
	if !ok {
		ginutils.RespError(c, "未找到用户信息")
	}

	if user, ok = v.(model.User); !ok {
		ginutils.RespError(c, "用户解析失败")
	}

	userService := userservice.NewUserService()
	if err := userService.UpdatePassword(user, req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	ginutils.RespOK(c, "用户密码修改成功")
}

// TODO: 将来需要修改这里的逻辑. 重置密码应该是发送修改密码链接到用户邮箱, 而不应该直接随机密码.
func (u User) ResetPassword(c *gin.Context) {
	var req dto.UserUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	userService := userservice.NewUserService()
	if err := userService.ResetPassword(req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	// c.SendOK("用户密码重置成功!")
	ginutils.RespOK(c, "用户密码重置链接已发送")
}

func (u User) Delete(c *gin.Context) {
	var req dto.UserDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	userService := userservice.NewUserService()
	if err := userService.Delete(req); err != nil {
		ginutils.RespError(c, err.Error())
		return
	}

	ginutils.RespOK(c, "用户删除成功")
}
