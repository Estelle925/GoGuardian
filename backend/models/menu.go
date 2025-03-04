package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title string `json:"title,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

// Value 实现 driver.Valuer 接口
func (m MenuMeta) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口
func (m *MenuMeta) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

// Menu 菜单模型
type Menu struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	ParentID          *int      `gorm:"default:0" json:"parent_id" example:"0"`
	Name              string    `gorm:"size:255;not null" json:"name" example:"系统管理"`
	Path              string    `gorm:"size:255;not null" json:"path" example:"/system"`
	Component         string    `gorm:"size:255;not null" json:"component" example:"@/views/system/index"`
	Icon              string    `gorm:"size:255" json:"icon,omitempty" example:"setting"`
	Order             int       `gorm:"default:0" json:"order" example:"1"`
	Meta              MenuMeta  `gorm:"type:json" json:"meta,omitempty"`
	IsVisible         bool      `gorm:"default:true" json:"is_visible" example:"true"`
	ButtonAssociation bool      `gorm:"default:false" json:"button_association" example:"false"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menu"
}