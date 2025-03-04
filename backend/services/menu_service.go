package services

import (
	"gorm.io/gorm"
	"tenant-center/models"
)

// MenuService 菜单服务
type MenuService struct {
	db *gorm.DB
}

// NewMenuService 创建菜单服务实例
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{db: db}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *models.Menu) error {
	return s.db.Create(menu).Error
}

// UpdateMenu 更新菜单信息
func (s *MenuService) UpdateMenu(menu *models.Menu) error {
	// 创建一个map来存储需要更新的字段
	updates := map[string]interface{}{
		"parent_id":           menu.ParentID,
		"name":               menu.Name,
		"path":               menu.Path,
		"component":          menu.Component,
		"icon":               menu.Icon,
		"order":              menu.Order,
		"meta":               menu.Meta,
		"is_visible":         menu.IsVisible,
		"button_association": menu.ButtonAssociation,
	}

	// 使用 Updates 方法只更新指定字段，让 GORM 自动处理时间戳
	return s.db.Model(menu).Updates(updates).Error
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id int) (*models.Menu, error) {
	var menu models.Menu
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// ListMenus 获取菜单列表
func (s *MenuService) ListMenus(page, pageSize int) ([]models.Menu, int64, error) {
	var menus []models.Menu
	var total int64

	// 计算总记录数
	if err := s.db.Model(&models.Menu{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// GetMenusByParentID 获取指定父级菜单的子菜单列表
func (s *MenuService) GetMenusByParentID(parentID int) ([]models.Menu, error) {
	var menus []models.Menu
	if err := s.db.Where("parent_id = ?", parentID).Order("order").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// BindMenuPermission 为菜单绑定权限
func (s *MenuService) BindMenuPermission(menuID int, permissionCode string, permissionName string) error {
	permission := &models.Permission{
		Code:   permissionCode,
		Name:   permissionName,
		Type:   "menu",
		MenuID: &menuID,
	}

	return s.db.Create(permission).Error
}

// GetMenuPermissions 获取菜单的权限列表
func (s *MenuService) GetMenuPermissions(menuID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Where("menu_id = ? AND type = 'menu'", menuID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
