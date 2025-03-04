package models

import (
	"time"
)

// Permission 权限模型
type Permission struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Code      string    `gorm:"size:255;not null;unique" json:"code" example:"user:create"`
	Name      string    `gorm:"size:255;not null" json:"name" example:"创建用户"`
	Type      string    `gorm:"type:enum('menu','button');not null" json:"type" example:"menu"`
	MenuID    *int      `gorm:"default:null" json:"menu_id,omitempty" example:"1"`
	ButtonID  *int      `gorm:"default:null" json:"button_id,omitempty" example:"1"`
	ParentID  *int      `gorm:"default:null" json:"parent_id,omitempty" example:"0"` // 父级权限ID
	Roles     []Role    `gorm:"many2many:role_permission;" json:"roles,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}
