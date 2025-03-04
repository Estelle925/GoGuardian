package services

import (
	"gorm.io/gorm"
	"tenant-center/models"
)

// RoleService 角色服务
type RoleService struct {
	db *gorm.DB
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *models.Role) error {
	return s.db.Create(role).Error
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(role *models.Role) error {
	// 创建一个map来存储需要更新的字段
	updates := map[string]interface{}{
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
	}

	// 使用 Updates 方法只更新指定字段，让 GORM 自动处理时间戳
	return s.db.Model(role).Updates(updates).Error
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// PageRoles 获取角色列表
func (s *RoleService) PageRoles(page, pageSize int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	// 计算总记录数
	if err := s.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// BindRolePermissionsByCode 通过权限编码绑定角色权限
func (s *RoleService) BindRolePermissionsByCode(roleID int, permissionCodes []string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除角色现有的所有权限
		if err := tx.Exec("DELETE FROM role_permission WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}

		// 查找所有指定编码的权限
		var permissions []models.Permission
		if err := tx.Where("id IN ?", permissionCodes).Find(&permissions).Error; err != nil {
			return err
		}

		// 添加新的权限关联
		for _, permission := range permissions {
			if err := tx.Exec("INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?)", roleID, permission.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *RoleService) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *RoleService) GetRolePermissions(roleID int) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := s.db.Table("permission").
		Select("permission.*").
		Joins("JOIN role_permission ON role_permission.permission_id = permission.id").
		Where("role_permission.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
