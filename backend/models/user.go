package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Username  string    `gorm:"size:255;not null;unique" json:"username" example:"admin"`
	Password  string    `gorm:"size:255;not null" json:"password,omitempty" example:"password123"`
	Roles     []Role    `gorm:"many2many:user_role;" json:"roles,omitempty"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// BeforeSave 保存前的钩子函数
func (u *User) BeforeSave(tx *gorm.DB) error {
	return nil
}
