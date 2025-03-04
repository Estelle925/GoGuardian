package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tenant-center/models"
	"tenant-center/services"
)

// @title 角色管理API
// @version 1.0
// @description 角色管理相关的API接口，包括创建角色、更新角色、获取角色列表、绑定权限和获取角色权限等功能

// RoleController 角色控制器
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{
		roleService: services.NewRoleService(db),
	}
}

// CreateRole @Summary 创建角色
// @Description 创建新角色，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body CreateRoleRequest true "角色信息"
// @Success 201 {object} CreateRoleResponse "角色创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "创建角色失败"
// @Security ApiKeyAuth
// @Router /api/roles [post]

// UpdateRole @Summary 更新角色信息
// @Description 更新指定角色的信息，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param role body UpdateRoleRequest true "角色信息"
// @Success 200 {object} UpdateRoleResponse "角色信息更新成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "更新角色信息失败"
// @Security ApiKeyAuth
// @Router /api/roles/{id} [put]

// PageRoles @Summary 获取角色列表
// @Description 获取所有角色的列表，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Success 200 {array} ListRolesResponse "角色列表"
// @Failure 500 {object} ErrorResponse "获取角色列表失败"
// @Security ApiKeyAuth
// @Router /api/roles [get]

// BindPermissions @Summary 为角色绑定权限
// @Description 为指定角色绑定一个或多个权限，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param permissionIDs body BindPermissionsRequest true "权限ID列表" example:[1,2,3]
// @Success 200 {object} BindPermissionsResponse "权限绑定成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "绑定权限失败"
// @Security ApiKeyAuth
// @Router /api/roles/{id}/permissions [post]

// PermissionTreeNode 权限树节点
type PermissionTreeNode struct {
	ID       int                  `json:"id"`
	Name     string               `json:"name"`
	Enable   bool                 `json:"enable"`
	Icon     string               `json:"icon,omitempty"`
	Children []PermissionTreeNode `json:"children"`
}

// GetDetail @Summary 获取角色详情
// @Description 获取指定角色的详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} models.Role "角色详情"
// @Failure 400 {object} ErrorResponse "无效的角色ID"
// @Failure 500 {object} ErrorResponse "获取角色详情失败"
// @Security ApiKeyAuth
// @Router /api/roles/detail/{id} [get]
func (c *RoleController) GetDetail(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	role, err := c.roleService.GetRoleByID(roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色详情失败"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

// GetRolePermissions @Summary 获取角色的权限列表
// @Description 获取指定角色的所有权限列表，按层级结构返回
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {array} PermissionTreeNode "权限树形列表"
// @Failure 400 {object} ErrorResponse "无效的角色ID"
// @Failure 500 {object} ErrorResponse "获取角色权限列表失败"
// @Security ApiKeyAuth
// @Router /api/roles/{id}/permissions [get]
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	// 获取所有权限
	allPermissions, err := c.roleService.GetAllPermissions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取所有权限列表失败"})
		return
	}

	// 获取角色已分配的权限
	rolePermissions, err := c.roleService.GetRolePermissions(roleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色权限列表失败"})
		return
	}

	// 创建角色权限ID集合，用于快速查找
	rolePermMap := make(map[int]bool)
	for _, perm := range rolePermissions {
		rolePermMap[perm.ID] = true
	}

	// 构建权限树
	permissionMap := make(map[int][]models.Permission)
	rootPermissions := make([]models.Permission, 0)

	// 按父ID对所有权限进行分组
	for _, perm := range allPermissions {
		if perm.ParentID == nil || *perm.ParentID == 0 {
			rootPermissions = append(rootPermissions, perm)
		} else {
			permissionMap[*perm.ParentID] = append(permissionMap[*perm.ParentID], perm)
		}
	}

	// 递归构建权限树
	buildPermissionTree := func(perms []models.Permission) []PermissionTreeNode {
		var buildTree func(permission models.Permission) PermissionTreeNode
		buildTree = func(permission models.Permission) PermissionTreeNode {
			node := PermissionTreeNode{
				ID:       permission.ID,
				Name:     permission.Name,
				Enable:   rolePermMap[permission.ID], // 根据角色权限设置启用状态
				Children: make([]PermissionTreeNode, 0),
			}

			if children, ok := permissionMap[permission.ID]; ok {
				for _, child := range children {
					node.Children = append(node.Children, buildTree(child))
				}
			}

			return node
		}

		treeNodes := make([]PermissionTreeNode, 0)
		for _, perm := range perms {
			treeNodes = append(treeNodes, buildTree(perm))
		}
		return treeNodes
	}

	// 生成最终的权限树
	permissionTree := buildPermissionTree(rootPermissions)

	ctx.JSON(http.StatusOK, permissionTree)
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var role struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	newRole := &models.Role{
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
	}

	if err := c.roleService.CreateRole(newRole); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "角色创建成功", "id": newRole.ID})
}

// UpdateRole @Summary 更新角色信息
// @Description 更新指定角色的信息，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param role body object true "角色信息"
// @Success 200 {object} object "角色信息更新成功"
// @Failure 400 {object} object "无效的请求参数"
// @Failure 500 {object} object "更新角色信息失败"
// @Security ApiKeyAuth
// @Router /api/roles/{id} [put]
// UpdateRole 更新角色信息
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var updateData struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	role := &models.Role{
		ID:          roleID,
		Name:        updateData.Name,
		Code:        updateData.Code,
		Description: updateData.Description,
	}

	if err := c.roleService.UpdateRole(role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新角色信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "角色信息更新成功"})
}

// GetRolesRequest 获取角色列表请求参数
type GetRolesRequest struct {
	Page     int `json:"page" example:"1" binding:"required"`
	PageSize int `json:"pageSize" example:"10" binding:"required"`
}

// GetRolesResponse 角色列表响应
type GetRolesResponse struct {
	Data     []models.Role `json:"data"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// PageRoles @Summary 获取角色列表
// @Description 获取所有角色的列表，支持分页，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body GetRolesRequest true "分页参数"
// @Success 200 {object} GetRolesResponse "角色列表"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "获取角色列表失败"
// @Security ApiKeyAuth
// @Router /api/roles [post]
// PageRoles 获取角色列表
func (c *RoleController) PageRoles(ctx *gin.Context) {
	var req GetRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	roles, total, err := c.roleService.PageRoles(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     roles,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// BindPermissions @Summary 为角色绑定权限
// @Description 为指定角色绑定一个或多个权限，需要管理员权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param permissions body BindPermissionsRequest true "权限编码列表"
// @Success 200 {object} object "权限绑定成功"
// @Failure 400 {object} object "无效的请求参数"
// @Failure 500 {object} object "绑定权限失败"
// @Security ApiKeyAuth
// @Router /api/roles/{id}/permissions [post]
// BindPermissions 为角色绑定权限
func (c *RoleController) BindPermissions(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := c.roleService.BindRolePermissionsByCode(roleID, req.Permissions); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "绑定权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true, "message": "权限绑定成功"})
}
