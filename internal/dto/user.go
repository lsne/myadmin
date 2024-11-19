/*
 * @Author: Liu Sainan
 * @Date: 2024-01-06 23:25:07
 */

package dto

import (
	"fmt"
	"myadmin/internal/model"
	"myadmin/internal/utils/htmlutils"
	"myadmin/internal/utils/validateutils"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserRegisterReq struct {
	Username string `json:"username" binding:"required" validate:"required,letterDigits,min=3,max=32"`
	Name     string `json:"name" binding:"required" validate:"required,min=6,max=32"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=100"`
	Phone    string `json:"phone" binding:"required" validate:"required,digits,min=7,max=16"`
	Email    string `json:"email" binding:"required" validate:"required,min=4,max=50"`
	Gender   uint   `json:"gender"`
	Role     int    `json:"role"`
}

func (r *UserRegisterReq) AvoidXSS() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Name = strings.TrimSpace(r.Name)
	r.Phone = strings.TrimSpace(r.Phone)
	r.Email = strings.TrimSpace(r.Email)

	if r.Username != htmlutils.AvoidXSS(r.Username) {
		return fmt.Errorf("Username(%s) 不符合规范", r.Username)
	}

	if r.Name != htmlutils.AvoidXSS(r.Name) {
		return fmt.Errorf("Name(%s) 不符合规范", r.Name)
	}

	if r.Email != htmlutils.AvoidXSS(r.Email) {
		return fmt.Errorf("Email(%s) 不符合规范", r.Email)
	}

	return nil
}

func (r *UserRegisterReq) Validator() error {

	validate := validator.New()

	if err := validate.RegisterValidation("letterDigits", validateutils.LetterDigits); err != nil {
		return err
	}

	if err := validate.RegisterValidation("digits", validateutils.Digits); err != nil {
		return err
	}
	return validate.Struct(r)
}

type UserListReq struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type UserListResp struct {
	Total int64        `json:"total"`
	Users []model.User `json:"users"`
}

// User 用户更新请求
type UserUpdateReq struct {
	ID        uint   `json:"id"`
	Username  string `json:"username" binding:"required" validate:"required,letterDigits,min=3,max=32"`
	Name      string `json:"name" binding:"required" validate:"required,min=6,max=32"`
	Phone     string `json:"phone" binding:"required" validate:"required,digits,min=7,max=16"`
	Email     string `json:"email" binding:"required" validate:"required,min=4,max=50"`
	Gender    uint   `json:"gender"`
	Status    int    `json:"status"`
	Signature string `json:"signature"`
	Introduce string `json:"introduce"`
}

func (r *UserUpdateReq) AvoidXSS() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Name = strings.TrimSpace(r.Name)
	r.Phone = strings.TrimSpace(r.Phone)
	r.Email = strings.TrimSpace(r.Email)

	if r.Username != htmlutils.AvoidXSS(r.Username) {
		return fmt.Errorf("Username(%s) 不符合规范", r.Username)
	}

	if r.Name != htmlutils.AvoidXSS(r.Name) {
		return fmt.Errorf("Name(%s) 不符合规范", r.Name)
	}

	if r.Email != htmlutils.AvoidXSS(r.Email) {
		return fmt.Errorf("Email(%s) 不符合规范", r.Email)
	}

	return nil
}

func (r *UserUpdateReq) Validator() error {

	validate := validator.New()

	if err := validate.RegisterValidation("letterDigits", validateutils.LetterDigits); err != nil {
		return err
	}

	if err := validate.RegisterValidation("digits", validateutils.Digits); err != nil {
		return err
	}
	return validate.Struct(r)
}

type UserUpdatePasswordReq struct {
	Password        string `json:"password" binding:"required,min=6,max=20"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=20"`
}

// User 用户删除请求
type UserDeleteReq struct {
	ID uint `json:"id"`
}
