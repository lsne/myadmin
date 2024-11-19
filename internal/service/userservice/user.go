/*
 * @Author: Liu Sainan
 * @Date: 2024-01-07 22:42:06
 */

package userservice

import (
	"fmt"
	"myadmin/internal/config"
	"myadmin/internal/dao"
	"myadmin/internal/dto"
	"myadmin/internal/model"
	"myadmin/internal/utils/mailutils"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) SendRegisterMail(user model.User) error {

	mail := mailutils.NewMail(
		config.GlobalConfig.Mail.Sender,
		[]string{user.Email},
		"text/html",
		config.GlobalConfig.Mail.Server,
		config.GlobalConfig.Mail.Port,
	)

	mail.Body = "<p><b>您好," + user.Name + ":</b></p>" +
		"<p>您在 " + config.GlobalConfig.Server.Name + " 的用户账户已经开通, 请尽快登录下面的链接并修改密码, 以激活帐号.</p>" +
		// "<a href=\"" + configs.GserConfig.Host + "\">" + configs.GserConfig.Host + "</a>" +
		"<p>账户: " + user.Username + "</p>" +
		"<p>初始密码同账号</p>" +
		"<p>" + config.GlobalConfig.Server.Name + "</p>" +
		"<p>(这是一封自动产生的 email, 请勿回复。)</p>"

	mail.SendMail("注册成功")
	return nil
}

func (u *UserService) RegisterPreCheck(user *model.User) error {
	dao := dao.NewUserDao()
	if _, err := dao.QueryFirstByUsername(user.Username); err == nil {
		return fmt.Errorf("Username(%s) 已经存在", user.Username)
	}
	if _, err := dao.QueryFirstByPhone(user.Phone); err == nil {
		return fmt.Errorf("Phone(%s) 已经存在", user.Phone)
	}
	if _, err := dao.QueryFirstByEmail(user.Email); err == nil {
		return fmt.Errorf("Email(%s) 已经存在", user.Email)
	}
	return nil
}

func (u *UserService) Register(req dto.UserRegisterReq) error {

	user := model.User{
		Username:  req.Username,
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Gender:    req.Gender,
		Role:      req.Role,
		AvatarURL: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
	}

	if err := u.RegisterPreCheck(&user); err != nil {
		return err
	}

	user.Password = user.EncryptPassword(req.Password, user.Salt())
	dao := dao.NewUserDao()
	if err := dao.Create(&user); err != nil {
		return err
	}

	return u.SendRegisterMail(user)
}

func (u *UserService) Login(req dto.UserLoginReq) (model.User, error) {

	dao := dao.NewUserDao()
	user, err := dao.QueryFirstByUsername(req.Username)
	if err != nil {
		return model.User{}, fmt.Errorf("Username(%s) 用户未注册: %v", req.Username, err)
	}

	if !user.CheckPassword(req.Password) {
		return model.User{}, fmt.Errorf("用户或密码错误")
	}

	if user.Status != 2 {
		return model.User{}, fmt.Errorf("用户未激活")
	}

	return user, nil
}

func (u *UserService) GetUserByID(id uint) (model.User, error) {

	dao := dao.NewUserDao()
	user, err := dao.QueryFirstByID(id)
	if err != nil {
		return model.User{}, fmt.Errorf("UserID(%d) 获取用户失败: %v", id, err)
	}

	return user, nil
}

func (u *UserService) GetSession(uid uint) (user model.User, err error) {
	dao := dao.NewUserDao()
	return dao.GetSession(uid)
}

func (u *UserService) SetSession(user *model.User) error {
	dao := dao.NewUserDao()
	return dao.SetSession(user)
}

func (u *UserService) DelSession(uid uint) error {
	dao := dao.NewUserDao()
	return dao.DelSession(uid)
}

func (u *UserService) GetRefreshToken(uid uint) (r string, err error) {
	dao := dao.NewUserDao()
	return dao.GetRefreshToken(uid)
}

func (u *UserService) DelRefreshToken(uid uint) error {
	dao := dao.NewUserDao()
	return dao.DelRefreshToken(uid)
}

func (u *UserService) Users(req dto.UserListReq) (resp dto.UserListResp, err error) {
	return dao.NewUserDao().QuerySearch(req)
}

func (u *UserService) Update(req dto.UserUpdateReq) error {
	dao := dao.NewUserDao()
	if err := dao.UpdateManyColumn(req); err != nil {
		return err
	}

	if req.Status != 0 {
		return u.DelSession(req.ID)
	}

	user, err := dao.QueryFirstByID(req.ID)
	if err != nil {
		return err
	}
	return u.SetSession(&user)
}

func (u *UserService) UpdatePassword(operator model.User, req dto.UserUpdatePasswordReq) (err error) {

	var user model.User

	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("二次校验密码不一致!")
	}

	dao := dao.NewUserDao()
	if user, err = dao.QueryFirstByID(operator.ID); err != nil {
		return err
	}

	if !user.CheckPassword(req.Password) {
		return fmt.Errorf("原始密码错误!")
	}

	user.Password = ""
	user.Password = user.EncryptPassword(req.NewPassword, user.Salt())

	if err := dao.UpdateByID(user.ID, "password", user.Password); err != nil {
		return err
	}
	u.DelRefreshToken(user.ID)
	return dao.DelSession(user.ID)
}

func (u *UserService) ResetPassword(req dto.UserUpdateReq) (err error) {

	var user model.User

	dao := dao.NewUserDao()
	if user, err = dao.QueryFirstByID(req.ID); err != nil {
		return err
	}

	password := fmt.Sprintf("%s_%s", req.Email, req.Phone)

	user.Password = ""
	user.Password = user.EncryptPassword(password, user.Salt())

	if err := dao.UpdateByID(user.ID, "password", user.Password); err != nil {
		return err
	}
	u.DelRefreshToken(req.ID)
	return dao.DelSession(user.ID)
}

func (u *UserService) Delete(req dto.UserDeleteReq) error {
	dao := dao.NewUserDao()
	if err := dao.DeleteByID(req.ID); err != nil {
		return err
	}

	u.DelRefreshToken(req.ID)
	return u.DelSession(req.ID)
}
