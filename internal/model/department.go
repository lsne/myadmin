/*
 * @Author: Liu Sainan
 * @Date: 2024-01-27 16:37:31
 */

package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// Department 部门表
type Department struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	ParentID  uint                  `gorm:"column:parent_id;not null;default:0;index;"`
	Ancestors string                `gorm:"column:ancestors;size:255;not null;default:'0.';index;"`
	Sort      uint                  `gorm:"column:sort;not null;default:0;index;"`
	Name      string                `gorm:"column:name;type:string;size:32;not null;default:'';index;comment:部门名称"`
	Status    int                   `gorm:"column:status;type:uint;size:8;not null;default:0;comment:未激活(0),冻结(1),激活(2)"` // 状态
	Memo      string                `gorm:"column:memo;type:string;size:255;not null;default:'';comment:备注"`              // 状态
	CreateBy  uint                  `gorm:"column:create_by;"`
	UpdateBy  uint                  `gorm:"column:update_by;"`
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}

// DepartmentUserRelation 部门用户关系表
type DepartmentUserRelation struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	DeptID    uint                  `gorm:"column:dept_id;index;not null;"`
	UserID    uint                  `gorm:"column:user_id;index;not null;"`
	RoleID    uint                  `gorm:"column:role_id;index;not null;default:0;"`
	Status    int                   `gorm:"column:status;type:uint;size:8;not null;default:0;comment:未激活(0),冻结(1),激活(2)"` // 状态
	CreateBy  uint                  `gorm:"column:create_by;"`
	UpdateBy  uint                  `gorm:"column:update_by;"`
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}
