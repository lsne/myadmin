/*
 * @Author: Liu Sainan
 * @Date: 2024-01-27 23:18:23
 */

package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

//  这个文件中的表暂时还没有用到, 用到的时候再整理

type WebMenu struct {
	ID            uint                                       `gorm:"column:id;autoIncrement;primary_key;not null;"`
	MenuLevel     uint                                       `json:"-"`
	ParentId      string                                     `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path          string                                     `json:"path" gorm:"comment:路由path"`        // 路由path
	Name          string                                     `json:"name" gorm:"comment:路由name"`        // 路由name
	Hidden        bool                                       `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	Component     string                                     `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Sort          int                                        `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	WebMenuMeta   `json:"meta" gorm:"embedded;comment:附加属性"` // 附加属性
	SysAuthoritys []SysAuthority                             `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []WebMenu                                  `json:"children" gorm:"-"`
	Parameters    []WebMenuParameter                         `json:"parameters"`
	// MenuBtn       []WebMenuBtn                               `json:"menuBtn"` // 这样搞得设置外键。 否则自动初始化表会失败
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}

type WebMenuMeta struct {
	ActiveName  string `json:"activeName" gorm:"comment:高亮菜单"`
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:是否缓存"`           // 是否缓存
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:是否是基础路由（开发中）"` // 是否是基础路由（开发中）
	Title       string `json:"title" gorm:"comment:菜单名"`                // 菜单名
	Icon        string `json:"icon" gorm:"comment:菜单图标"`                // 菜单图标
	CloseTab    bool   `json:"closeTab" gorm:"comment:自动关闭tab"`         // 自动关闭tab
}

type WebMenuParameter struct {
	ID        uint `gorm:"column:id;autoIncrement;primary_key;not null;"`
	WebMenuID uint
	Type      string                `json:"type" gorm:"comment:地址栏携带参数为params还是query"`                     // 地址栏携带参数为params还是query
	Key       string                `json:"key" gorm:"comment:地址栏携带参数的key"`                                // 地址栏携带参数的key
	Value     string                `json:"value" gorm:"comment:地址栏携带参数的值"`                                // 地址栏携带参数的值
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}

type WebMenuWith struct {
	WebMenu
	MenuId      string        `json:"menuId" gorm:"comment:菜单ID"`
	AuthorityId uint          `json:"-" gorm:"comment:角色ID"`
	Children    []WebMenuWith `json:"children" gorm:"-"`
	// Parameters  []WebMenuParameter `json:"parameters"`  // 这样搞得设置外键。 否则自动初始化表会失败
	Btns map[string]uint `json:"btns" gorm:"-"`
}

// type SysAuthorityMenu struct {
// 	MenuId      string `json:"menuId" gorm:"comment:菜单ID;"`
// 	AuthorityId string `json:"-" gorm:"comment:角色ID;"`
// }

// func (s SysAuthorityMenu) TableName() string {
// 	return "sys_authority_menus"
// }

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []WebMenu       `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []User          `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}

type WebMenuBtn struct {
	ID            uint      `gorm:"column:id;autoIncrement;primary_key;not null;"`
	Name          string    `json:"name" gorm:"comment:按钮关键key"`
	Desc          string    `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID uint      `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at"`
}
