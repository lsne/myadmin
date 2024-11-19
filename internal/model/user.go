/*
 * @Author: Liu Sainan
 * @Date: 2024-01-02 01:29:46
 */

package model

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"myadmin/internal/config"
	"strconv"
	"time"

	"gorm.io/plugin/soft_delete"
)

// Role 角色表, 暂时只统一规定角色, 同一角色在每个部门, 部门组, 项目, 项目组中的权限相同
// 0 - 普通成员
// 1 - 部门经理
// 2 - 部门助理
// 3 - 项目经理
// 4 - 项目助理
// 5 - 管理员
// 6 - 超级管理员
type Role struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	DeptID    uint                  `gorm:"column:dept_id;autoIncrement;index;not null;"`
	Name      string                `gorm:"column:name;type:string;size:32;not null;default:'';index;comment:角色名称"`
	Status    int                   `gorm:"column:status;type:uint;size:8;not null;default:0;comment:未激活(0),冻结(1),激活(2)"` // 状态
	CreateBy  uint                  `gorm:"column:create_by;"`
	UpdateBy  uint                  `gorm:"column:update_by;"`
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}

// User 用户表
type User struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;" json:"id"`
	Username  string                `gorm:"column:username;type:string;size:32;not null;default:'';unique;comment:用户名" json:"username"`
	Name      string                `gorm:"column:name;type:string;size:32;not null;default:'';index;comment:用户姓名" json:"name"`
	Password  string                `gorm:"column:password;type:string;size:80;not null;default:'';comment:用户密码" json:"-"`
	Phone     string                `gorm:"column:phone;type:string;size:16;not null;default:'';unique;comment:手机号" json:"phone"`
	Email     string                `gorm:"column:email;type:string;size:50;not null;default:'';unique;comment:邮箱" json:"email"`
	Gender    uint                  `gorm:"column:gender;type:uint;size:8;not null;default:0;comment:性别" json:"gender"`
	AvatarURL string                `gorm:"column:avatar_url;type:string;size:128;not null;default:'';comment:头像" json:"avatar_url"`                 // 头像
	Role      int                   `gorm:"column:role;type:int;size:16;not null;default:0;comment:普通用户(0),管理员(10),高级管理员(20),超级管理员(30)" json:"role"` // 角色, 决定允许对系统的用户端访问还是管理端访问, 与 Role 表没有关系
	Status    int                   `gorm:"column:status;type:uint;size:8;not null;default:0;comment:未激活(0),冻结(1),激活(2)" json:"status"`              // 状态. 目前还没用上. 只要注册就能登录
	Signature string                `gorm:"column:signature;type:string;size:128" json:"signature"`                                                  // 个人签名
	Introduce string                `gorm:"column:introduce;type:string;size:512" json:"introduce"`
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt" json:"-"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at" json:"update_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// CheckPassword 验证密码是否正确
func (user *User) CheckPassword(password string) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return user.EncryptPassword(password, user.Salt()) == user.Password
}

// Salt 每个用户都有一个不同的盐
func (user *User) Salt() string {
	var userSalt string
	if user.Password == "" {
		userSalt = strconv.Itoa(int(time.Now().Unix()))
	} else {
		userSalt = user.Password[0:10]
	}
	return userSalt
}

// EncryptPassword 给密码加密
func (user *User) EncryptPassword(password, salt string) string {
	password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	hash := salt + password + config.GlobalConfig.Auth.PasswordSalt
	return salt + fmt.Sprintf("%x", md5.Sum([]byte(hash)))
}

const (
	// UserRoleNormal 普通用户
	UserRoleNormal = 0

	// UserRoleAdmin 管理员
	UserRoleAdmin = 10

	// UserRoleEditor 高级管理员
	UserRoleSeniorAdmin = 20

	// UserRoleRoot 超级管理员
	UserRoleRoot = 30
)

const (
	// UserStatusInActive 未激活
	UserStatusInActive = 0

	// UserStatusActived 已冻结
	UserStatusActived = 1

	// UserStatusFrozen 已激活
	UserStatusFrozen = 2
)

const (
	// UserSexMale 男
	UserGenderMale = 0

	// UserSexFemale 女
	UserGenderFemale = 1

	// MaxUserNameLen 用户名的最大长度
	MaxUserNameLen = 20

	// MinUserNameLen 用户名的最小长度
	MinUserNameLen = 4

	// MaxPassLen 密码的最大长度
	MaxPassLen = 20

	// MinPassLen 密码的最小长度
	MinPassLen = 6

	// MaxSignatureLen 个性签名最大长度
	MaxSignatureLen = 200

	// MaxLocationLen 居住地的最大长度
	MaxLocationLen = 200

	// MaxIntroduceLen 个人简介的最大长度
	MaxIntroduceLen = 500
)
