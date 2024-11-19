/*
 * @Author: Liu Sainan
 * @Date: 2024-01-27 23:40:52
 */

package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

//  这个文件中的表暂时还没有用到, 用到的时候再整理

type ServerApi struct {
	ID          uint                  `gorm:"column:id;autoIncrement;primary_key;not null;"`
	Path        string                `json:"path" gorm:"comment:api路径"`                                     // api路径
	Description string                `json:"description" gorm:"comment:api中文描述"`                            // api中文描述
	ApiGroup    string                `json:"apiGroup" gorm:"comment:api组"`                                  // api组
	Method      string                `json:"method" gorm:"default:POST;comment:方法"`                         // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	IsDel       soft_delete.DeletedAt `gorm:"column:is_del;size:8;softDelete:flag,DeletedAtField:DeletedAt"` // use `1` `0`
	CreatedAt   time.Time             `gorm:"column:created_at"`
	UpdatedAt   time.Time             `gorm:"column:updated_at"`
	DeletedAt   time.Time             `gorm:"column:deleted_at"`
}

func (ServerApi) TableName() string {
	return "server_apis"
}
