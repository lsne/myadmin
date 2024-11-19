/*
 * @Author: Liu Sainan
 * @Date: 2024-01-27 22:56:49
 */

package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// project 项目表
type Project struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	DeptID    uint                  `gorm:"column:id;index;not null;"`
	Name      string                `gorm:"column:name;type:string;size:32;not null;default:'';index;comment:项目名称"`
	Type      int                   `gorm:"column:type;type:uint;size:8;not null;default:0;comment:研发(0),预研(1),运维(2),运营(3)"` // 状态
	Status    int                   `gorm:"column:status;type:uint;size:8;not null;default:0;comment:未激活(0),冻结(1),激活(2)"`    // 状态
	Memo      string                `gorm:"column:memo;type:string;size:255;not null;default:'';comment:备注"`                 // 状态
	CreateBy  uint                  `gorm:"column:create_by;"`
	UpdateBy  uint                  `gorm:"column:update_by;"`
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt time.Time             `gorm:"column:created_at"`
	UpdatedAt time.Time             `gorm:"column:updated_at"`
	DeletedAt time.Time             `gorm:"column:deleted_at"`
}

// ProjectUserRelation 部门用户关系表
type ProjectUserRelation struct {
	ID        uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	ProjectID uint                  `gorm:"column:project_id;index;not null;"`
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
