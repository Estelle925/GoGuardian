package services

import (
	"gorm.io/gorm"
	"tenant-center/models"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(permission *models.Permission) error {
	return s.db.Create(permission).Error
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(permission *models.Permission) error {
	// 创建一个map来存储需要更新的字段
	updates := map[string]interface{}{
		"code":      permission.Code,
		"name":      permission.Name,
		"type":      permission.Type,
		"menu_id":   permission.MenuID,
		"button_id": permission.ButtonID,
	}

	// 使用 Updates 方法只更新指定字段，让 GORM 自动处理时间戳
	return s.db.Model(permission).Updates(updates).Error
}

// GetPermissionByID 根据ID获取权限
func (s *PermissionService) GetPermissionByID(id int) (*models.Permission, error) {
	var permission models.Permission
	if err := s.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// PagePermissions 获取权限列表
func (s *PermissionService) PagePermissions(page, pageSize int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	// 计算总记录数
	if err := s.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetPermissionsByType 根据类型获取权限列表
func (s *PermissionService) GetPermissionsByType(permissionType string) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Where("type = ?", permissionType).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionsByMenuID 获取指定菜单的权限列表
func (s *PermissionService) GetPermissionsByMenuID(menuID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Where("menu_id = ?", menuID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionsByButtonID 获取指定按钮的权限列表
func (s *PermissionService) GetPermissionsByButtonID(buttonID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Where("button_id = ?", buttonID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
