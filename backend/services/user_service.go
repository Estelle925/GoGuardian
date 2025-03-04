package services

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"tenant-center/models"
	"time"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *models.User) error {
	// 将密码进行了Base64编码，直接存储
	user.Password = base64.StdEncoding.EncodeToString([]byte(user.Password))
	return s.db.Create(user).Error
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *models.User) error {
	// 创建一个map来存储需要更新的字段
	updates := map[string]interface{}{
		"username": user.Username,
	}

	// 只有当密码不为空时才更新密码
	if user.Password != "" {
		updates["password"] = base64.StdEncoding.EncodeToString([]byte(user.Password))
	}

	// 使用 Updates 方法只更新指定字段，让 GORM 自动处理时间戳
	return s.db.Model(user).Updates(updates).Error
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	// 查找用户
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("用户不存在")
	}

	// 直接比较Base64编码后的密码
	if user.Password != password {
		return "", errors.New("密码错误")
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte("your-secret-key")) // 注意：在实际应用中，密钥应该从配置文件中读取
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// BindUserRoles 为用户绑定角色
func (s *UserService) BindUserRoles(userID int, roleIDs []int) error {
	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先删除用户现有的所有角色
		if err := tx.Exec("DELETE FROM user_role WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		for _, roleID := range roleIDs {
			if err := tx.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// RouteItem 路由项结构
type RouteItem struct {
	Component string      `json:"component"`
	Meta      RouteMeta   `json:"meta"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Children  []RouteItem `json:"children,omitempty"`
}

// RouteMeta 路由元数据
type RouteMeta struct {
	Icon       string `json:"icon,omitempty"`
	DarkIcon   string `json:"darkIcon,omitempty"`
	ActiveIcon string `json:"activeIcon,omitempty"`
	Order      int    `json:"order,omitempty"`
	Title      string `json:"title"`
	Authority  []int  `json:"authority,omitempty"`
}

// PageUsers 获取用户列表，支持分页
func (s *UserService) PageUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 计算总记录数
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，排除密码字段
	offset := (page - 1) * pageSize
	if err := s.db.Select("id, username, created_at, updated_at").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserRoutes 获取用户的路由数据
func (s *UserService) GetUserRoutes(userID int) ([]RouteItem, error) {
	// 获取用户的角色
	var user models.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 获取所有菜单
	var menus []models.Menu
	if err := s.db.Order("parent_id, `order`").Find(&menus).Error; err != nil {
		return nil, err
	}

	// 构建菜单树
	return s.buildRouteTree(menus, 0), nil
}

// buildRouteTree 构建路由树
func (s *UserService) buildRouteTree(menus []models.Menu, parentID int) []RouteItem {
	var routes []RouteItem

	for _, menu := range menus {
		if menu.ParentID == nil && parentID == 0 || (menu.ParentID != nil && *menu.ParentID == parentID) {
			route := RouteItem{
				Component: menu.Component,
				Name:      menu.Name,
				Path:      menu.Path,
				Meta: RouteMeta{
					Title:      menu.Meta.Title,
					Icon:       menu.Icon,
					DarkIcon:   menu.Icon, // 可以根据需要设置不同的图标
					ActiveIcon: menu.Icon,
					Order:      menu.Order,
					Authority:  []int{1}, // 这里可以根据实际权限设置
				},
			}

			// 递归获取子菜单
			children := s.buildRouteTree(menus, menu.ID)
			if len(children) > 0 {
				route.Children = children
			}

			routes = append(routes, route)
		}
	}

	return routes
}
