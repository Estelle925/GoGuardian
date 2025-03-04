package models

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID          int          `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name        string       `gorm:"size:255;not null" json:"name" example:"管理员"`
	Code        string       `gorm:"size:255;not null;unique" json:"code" example:"ROLE_ADMIN"`
	Description string       `gorm:"type:text" json:"description" example:"系统管理员角色"`
	Permissions []Permission `gorm:"many2many:role_permission;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_role;" json:"users,omitempty"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}