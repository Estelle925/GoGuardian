package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tenant-center/models"
	"tenant-center/services"
)

// @title 权限管理API
// @version 1.0
// @description 权限管理相关的API接口，包括创建权限、更新权限、获取权限列表和根据类型获取权限等功能

// PermissionController 权限控制器
type PermissionController struct {
	permissionService *services.PermissionService
}

// NewPermissionController 创建权限控制器实例
func NewPermissionController(db *gorm.DB) *PermissionController {
	return &PermissionController{
		permissionService: services.NewPermissionService(db),
	}
}

// CreatePermissionRequest 创建权限请求参数
type CreatePermissionRequest struct {
	Code     string `json:"code" binding:"required" example:"user:create"` // 权限编码
	Name     string `json:"name" binding:"required" example:"创建用户"`        // 权限名称
	Type     string `json:"type" binding:"required" example:"menu"`        // 权限类型
	MenuID   *int   `json:"menu_id" example:"1"`                           // 菜单ID
	ButtonID *int   `json:"button_id" example:"1"`                         // 按钮ID
}

// CreatePermissionResponse 创建权限响应
type CreatePermissionResponse struct {
	Message string `json:"message" example:"权限创建成功"` // 响应消息
	ID      int    `json:"id" example:"1"`           // 权限ID
}

// UpdatePermissionRequest 更新权限请求参数
type UpdatePermissionRequest struct {
	Code     string `json:"code" example:"user:create"` // 权限编码
	Name     string `json:"name" example:"创建用户"`        // 权限名称
	Type     string `json:"type" example:"menu"`        // 权限类型
	MenuID   *int   `json:"menu_id" example:"1"`        // 菜单ID
	ButtonID *int   `json:"button_id" example:"1"`      // 按钮ID
}

// UpdatePermissionResponse 更新权限响应
type UpdatePermissionResponse struct {
	Message string `json:"message" example:"权限信息更新成功"` // 响应消息
}

// GetPermissionsByTypeResponse 根据类型获取权限列表响应
type GetPermissionsByTypeResponse []models.Permission

// CreatePermission @Summary 创建权限
// @Description 创建新权限，需要管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param permission body CreatePermissionRequest true "权限信息"
// @Success 201 {object} CreatePermissionResponse "权限创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "创建权限失败"
// @Security ApiKeyAuth
// @Router /api/permissions [post]
// CreatePermission 创建权限
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	// 定义请求参数结构体
	var permission struct {
		PermissionCode string `json:"permissionCode" binding:"required"` // 改正字段名为 permissionCode
		Name           string `json:"name" binding:"required"`           // 保持一致的命名
		Type           string `json:"type" binding:"required"`           // 保持一致的命名
		MenuID         *int   `json:"menuId"`                            // MenuId 允许为空
		ButtonID       *int   `json:"buttonId"`                          // ButtonId 允许为空
	}

	// 绑定请求的 JSON 数据到 permission 结构体
	if err := ctx.ShouldBindJSON(&permission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 创建新权限对象
	newPermission := &models.Permission{
		Code:     permission.PermissionCode, // 绑定对应的字段
		Name:     permission.Name,
		Type:     permission.Type,
		MenuID:   permission.MenuID,
		ButtonID: permission.ButtonID,
	}

	// 调用 service 层进行权限创建
	if err := c.permissionService.CreatePermission(newPermission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建权限失败"})
		return
	}

	// 返回创建成功的响应
	ctx.JSON(http.StatusCreated, gin.H{"message": "权限创建成功", "id": newPermission.ID})
}

// UpdatePermission @Summary 更新权限信息
// @Description 更新指定权限的信息，需要管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param permission body UpdatePermissionRequest true "权限信息"
// @Success 200 {object} UpdatePermissionResponse "权限信息更新成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "更新权限信息失败"
// @Security ApiKeyAuth
// @Router /api/permissions/{id} [put]
// UpdatePermission 更新权限信息
func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	permissionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	var updateData struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		MenuID   *int   `json:"menu_id"`
		ButtonID *int   `json:"button_id"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	permission := &models.Permission{
		ID:       permissionID,
		Code:     updateData.Code,
		Name:     updateData.Name,
		Type:     updateData.Type,
		MenuID:   updateData.MenuID,
		ButtonID: updateData.ButtonID,
	}

	if err := c.permissionService.UpdatePermission(permission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新权限信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "权限信息更新成功"})
}

// GetPermissionsRequest 获取权限列表请求参数
type GetPermissionsRequest struct {
	Page     int `json:"page" example:"1" binding:"required"`
	PageSize int `json:"pageSize" example:"10" binding:"required"`
}

// GetPermissionsResponse 权限列表响应
type GetPermissionsResponse struct {
	Data     []models.Permission `json:"data"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

// PagePermissions @Summary 获取权限列表
// @Description 获取所有权限的列表，支持分页，需要管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body GetPermissionsRequest true "分页参数"
// @Success 200 {object} GetPermissionsResponse "权限列表"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "获取权限列表失败"
// @Security ApiKeyAuth
// @Router /api/permissions [post]
// PagePermissions 获取权限列表
func (c *PermissionController) PagePermissions(ctx *gin.Context) {
	var req GetPermissionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	permissions, total, err := c.permissionService.PagePermissions(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     permissions,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// GetPermissionDetail @Summary 获取权限详情
// @Description 获取指定权限的详细信息，包括基本信息、关联的菜单和按钮信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} models.Permission "权限详情"
// @Failure 400 {object} ErrorResponse "无效的权限ID"
// @Failure 404 {object} ErrorResponse "权限不存在"
// @Failure 500 {object} ErrorResponse "获取权限详情失败"
// @Security ApiKeyAuth
// @Router /api/permissions/{id} [get]
func (c *PermissionController) GetPermissionDetail(ctx *gin.Context) {
	permissionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	permission, err := c.permissionService.GetPermissionByID(permissionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限详情失败"})
		return
	}

	ctx.JSON(http.StatusOK, permission)
}

// GetPermissionsByType @Summary 根据类型获取权限列表
// @Description 获取指定类型的所有权限列表，需要管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param type path string true "权限类型"
// @Success 200 {array} GetPermissionsByTypeResponse "权限列表"
// @Failure 500 {object} ErrorResponse "获取权限列表失败"
// @Security ApiKeyAuth
// @Router /api/permissions/type/{type} [get]
// GetPermissionsByType 根据类型获取权限列表
func (c *PermissionController) GetPermissionsByType(ctx *gin.Context) {
	permissionType := ctx.Param("type")

	permissions, err := c.permissionService.GetPermissionsByType(permissionType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, permissions)
}
