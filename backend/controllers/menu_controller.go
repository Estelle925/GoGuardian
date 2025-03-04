package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tenant-center/models"
	"tenant-center/services"
)

// @title 菜单管理API
// @version 1.0
// @description 菜单管理相关的API接口，包括创建菜单、更新菜单、获取菜单列表、获取子菜单、绑定权限和获取菜单权限等功能

// MenuController 菜单控制器
type MenuController struct {
	menuService *services.MenuService
}

// NewMenuController 创建菜单控制器实例
func NewMenuController(db *gorm.DB) *MenuController {
	return &MenuController{
		menuService: services.NewMenuService(db),
	}
}

// CreateMenuRequest 创建菜单请求参数
type CreateMenuRequest struct {
	Name     string `json:"name" binding:"required" example:"系统管理"`    // 菜单名称
	Path     string `json:"path" binding:"required" example:"/system"` // 菜单路径
	Icon     string `json:"icon" example:"setting"`                    // 菜单图标
	ParentID *int   `json:"parent_id" example:"0"`                     // 父级菜单ID
	Order    int    `json:"order" example:"1"`                         // 排序序号
}

// CreateMenuResponse 创建菜单响应
type CreateMenuResponse struct {
	Message string `json:"message" example:"菜单创建成功"` // 响应消息
	ID      int    `json:"id" example:"1"`           // 菜单ID
}

// UpdateMenuRequest 更新菜单请求参数
type UpdateMenuRequest struct {
	Name     string `json:"name" example:"用户管理"`   // 菜单名称
	Path     string `json:"path" example:"/user"`  // 菜单路径
	Icon     string `json:"icon" example:"user"`   // 菜单图标
	ParentID *int   `json:"parent_id" example:"1"` // 父级菜单ID
	Order    int    `json:"order" example:"2"`     // 排序序号
}

// UpdateMenuResponse 更新菜单响应
type UpdateMenuResponse struct {
	Message string `json:"message" example:"菜单信息更新成功"` // 响应消息
}

// GetMenusByParentIDResponse 获取子菜单列表响应
type GetMenusByParentIDResponse []models.Menu

// GetDetail @Summary 获取菜单详情
// @Description 获取指定菜单的详细信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} models.Menu "菜单详情"
// @Failure 400 {object} ErrorResponse "无效的菜单ID"
// @Failure 500 {object} ErrorResponse "获取菜单详情失败"
// @Security ApiKeyAuth
// @Router /api/menus/detail/{id} [get]
func (c *MenuController) GetDetail(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	menu, err := c.menuService.GetMenuByID(menuID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单详情失败"})
		return
	}

	ctx.JSON(http.StatusOK, menu)
}

// BindMenuPermissionRequest 绑定菜单权限请求参数
type BindMenuPermissionRequest struct {
	PermissionCode string `json:"permission_code" binding:"required" example:"menu:view"` // 权限编码
	Name           string `json:"name" binding:"required" example:"查看菜单"`                 // 权限名称
}

// BindMenuPermissionResponse 绑定菜单权限响应
type BindMenuPermissionResponse struct {
	Message string `json:"message" example:"权限绑定成功"` // 响应消息
}

// GetMenuPermissionsResponse 获取菜单权限列表响应
type GetMenuPermissionsResponse []models.Permission

// CreateMenu @Summary 创建菜单
// @Description 创建新菜单，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body CreateMenuRequest true "菜单信息"
// @Success 201 {object} CreateMenuResponse "菜单创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "创建菜单失败"
// @Security ApiKeyAuth
// @Router /api/menus [post]
// CreateMenu 创建菜单
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var menu struct {
		Name     string `json:"name" binding:"required"`
		Path     string `json:"path" binding:"required"`
		Icon     string `json:"icon"`
		ParentID *int   `json:"parent_id"`
		Order    int    `json:"order"`
	}

	if err := ctx.ShouldBindJSON(&menu); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	newMenu := &models.Menu{
		Name:     menu.Name,
		Path:     menu.Path,
		Icon:     menu.Icon,
		ParentID: menu.ParentID,
		Order:    menu.Order,
	}

	if err := c.menuService.CreateMenu(newMenu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜单失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "菜单创建成功", "id": newMenu.ID})
}

// UpdateMenu @Summary 更新菜单信息
// @Description 更新指定菜单的信息，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param menu body UpdateMenuRequest true "菜单信息"
// @Success 200 {object} UpdateMenuResponse "菜单信息更新成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "更新菜单信息失败"
// @Security ApiKeyAuth
// @Router /api/menus/{id} [put]
// UpdateMenu 更新菜单信息
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	var updateData struct {
		Name     string `json:"name"`
		Path     string `json:"path"`
		Icon     string `json:"icon"`
		ParentID *int   `json:"parent_id"`
		Order    int    `json:"order"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	menu := &models.Menu{
		ID:       menuID,
		Name:     updateData.Name,
		Path:     updateData.Path,
		Icon:     updateData.Icon,
		ParentID: updateData.ParentID,
		Order:    updateData.Order,
	}

	if err := c.menuService.UpdateMenu(menu); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新菜单信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "菜单信息更新成功"})
}

// GetMenusRequest 获取菜单列表请求参数
type GetMenusRequest struct {
	Page     int `json:"page" example:"1" binding:"required"`
	PageSize int `json:"pageSize" example:"10" binding:"required"`
}

// GetMenusResponse 菜单列表响应
type GetMenusResponse struct {
	Data     []models.Menu `json:"data"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

// ListMenus @Summary 获取菜单列表
// @Description 获取所有菜单的列表，支持分页，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body GetMenusRequest true "分页参数"
// @Success 200 {object} GetMenusResponse "菜单列表"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "获取菜单列表失败"
// @Security ApiKeyAuth
// @Router /api/menus [post]
// ListMenus 获取菜单列表
func (c *MenuController) ListMenus(ctx *gin.Context) {
	var req GetMenusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	menus, total, err := c.menuService.ListMenus(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     menus,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// GetMenusByParentID @Summary 获取子菜单列表
// @Description 获取指定父级菜单的子菜单列表，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param parentId path int true "父级菜单ID"
// @Success 200 {array} GetMenusByParentIDResponse "子菜单列表"
// @Failure 400 {object} ErrorResponse "无效的父级菜单ID"
// @Failure 500 {object} ErrorResponse "获取子菜单列表失败"
// @Security ApiKeyAuth
// @Router /api/menus/parent/{parentId} [get]
// GetMenusByParentID 获取指定父级菜单的子菜单列表
func (c *MenuController) GetMenusByParentID(ctx *gin.Context) {
	parentID, err := strconv.Atoi(ctx.Param("parentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的父级菜单ID"})
		return
	}

	menus, err := c.menuService.GetMenusByParentID(parentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取子菜单列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, menus)
}

// BindPermission @Summary 为菜单绑定权限
// @Description 为指定菜单绑定权限，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param permission body BindMenuPermissionRequest true "权限信息"
// @Success 200 {object} BindMenuPermissionResponse "权限绑定成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "绑定权限失败"
// @Security ApiKeyAuth
// @Router /api/menus/{id}/permission [post]
// BindPermission 为菜单绑定权限
func (c *MenuController) BindPermission(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	var permission struct {
		PermissionCode string `json:"permission_code" binding:"required"`
		Name           string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := c.menuService.BindMenuPermission(menuID, permission.PermissionCode, permission.Name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "绑定权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "权限绑定成功"})
}

// GetMenuPermissions @Summary 获取菜单权限列表
// @Description 获取指定菜单的权限列表，需要管理员权限
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {array} GetMenuPermissionsResponse "权限列表"
// @Failure 400 {object} ErrorResponse "无效的菜单ID"
// @Failure 500 {object} ErrorResponse "获取菜单权限列表失败"
// @Security ApiKeyAuth
// @Router /api/menus/{id}/permissions [get]
// GetMenuPermissions 获取菜单的权限列表
func (c *MenuController) GetMenuPermissions(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	permissions, err := c.menuService.GetMenuPermissions(menuID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜单权限列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, permissions)
}
