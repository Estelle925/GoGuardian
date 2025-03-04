package models

import (
	"time"
)

// Button 按钮模型
type Button struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name           string    `gorm:"size:255;not null" json:"name" example:"创建用户"`
	Action         string    `gorm:"size:255;not null" json:"action" example:"create"`
	MenuID         int       `gorm:"not null" json:"menu_id" example:"1"`
	PermissionCode string    `gorm:"size:255;not null" json:"permission_code" example:"user:create"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Button) TableName() string {
	return "button"
}
