package services

import (
	"gorm.io/gorm"
	"tenant-center/models"
)

// ButtonService 按钮服务
type ButtonService struct {
	db *gorm.DB
}

// NewButtonService 创建按钮服务实例
func NewButtonService(db *gorm.DB) *ButtonService {
	return &ButtonService{db: db}
}

// CreateButton 创建按钮
func (s *ButtonService) CreateButton(button *models.Button) error {
	return s.db.Create(button).Error
}

// UpdateButton 更新按钮信息
func (s *ButtonService) UpdateButton(button *models.Button) error {
	// 创建一个map来存储需要更新的字段
	updates := map[string]interface{}{
		"name":            button.Name,
		"permission_code": button.PermissionCode,
		"menu_id":         button.MenuID,
	}

	// 使用 Updates 方法只更新指定字段，让 GORM 自动处理时间戳
	return s.db.Model(button).Updates(updates).Error
}

// GetButtonByID 根据ID获取按钮
func (s *ButtonService) GetButtonByID(id int) (*models.Button, error) {
	var button models.Button
	if err := s.db.First(&button, id).Error; err != nil {
		return nil, err
	}
	return &button, nil
}

// ListButtons 获取按钮列表
func (s *ButtonService) ListButtons(page, pageSize int) ([]models.Button, int64, error) {
	var buttons []models.Button
	var total int64

	// 计算总记录数
	if err := s.db.Model(&models.Button{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&buttons).Error; err != nil {
		return nil, 0, err
	}

	return buttons, total, nil
}

// GetButtonsByMenuID 获取指定菜单的按钮列表
func (s *ButtonService) GetButtonsByMenuID(menuID int) ([]models.Button, error) {
	var buttons []models.Button
	if err := s.db.Where("menu_id = ?", menuID).Find(&buttons).Error; err != nil {
		return nil, err
	}
	return buttons, nil
}

// BindButtonPermission 为按钮绑定权限
func (s *ButtonService) BindButtonPermission(buttonID int, permissionCode string, permissionName string) error {
	permission := &models.Permission{
		Code:     permissionCode,
		Name:     permissionName,
		Type:     "button",
		ButtonID: &buttonID,
	}

	return s.db.Create(permission).Error
}

// GetButtonPermissions 获取按钮的权限列表
func (s *ButtonService) GetButtonPermissions(buttonID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Where("button_id = ? AND type = 'button'", buttonID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
