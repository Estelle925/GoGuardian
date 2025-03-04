package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tenant-center/models"
	"tenant-center/services"
)

// @title 按钮管理API
// @version 1.0
// @description 按钮管理相关的API接口，包括创建按钮、更新按钮、获取按钮列表、绑定权限和获取按钮权限等功能

// ButtonController 按钮控制器
type ButtonController struct {
	buttonService *services.ButtonService
}

// NewButtonController 创建按钮控制器实例
func NewButtonController(db *gorm.DB) *ButtonController {
	return &ButtonController{
		buttonService: services.NewButtonService(db),
	}
}

// CreateButtonRequest 创建按钮请求参数
type CreateButtonRequest struct {
	Name   string `json:"name" binding:"required" example:"新增"`   // 按钮名称
	Code   string `json:"code" binding:"required" example:"add"`  // 按钮编码
	MenuID int    `json:"menu_id" binding:"required" example:"1"` // 所属菜单ID
}

// CreateButtonResponse 创建按钮响应
type CreateButtonResponse struct {
	Message string `json:"message" example:"按钮创建成功"` // 响应消息
	ID      int    `json:"id" example:"1"`           // 按钮ID
}

// UpdateButtonRequest 更新按钮请求参数
type UpdateButtonRequest struct {
	Name   string `json:"name" example:"编辑"`   // 按钮名称
	Code   string `json:"code" example:"edit"` // 按钮编码
	MenuID int    `json:"menu_id" example:"1"` // 所属菜单ID
}

// UpdateButtonResponse 更新按钮响应
type UpdateButtonResponse struct {
	Message string `json:"message" example:"按钮信息更新成功"` // 响应消息
}

// GetButtonsByMenuIDResponse 获取菜单按钮列表响应
type GetButtonsByMenuIDResponse []models.Button

// BindButtonPermissionRequest 绑定按钮权限请求参数
type BindButtonPermissionRequest struct {
	PermissionCode string `json:"permission_code" binding:"required" example:"button:add"` // 权限编码
	Name           string `json:"name" binding:"required" example:"新增按钮"`                  // 权限名称
}

// BindButtonPermissionResponse 绑定按钮权限响应
type BindButtonPermissionResponse struct {
	Message string `json:"message" example:"权限绑定成功"` // 响应消息
}

// GetDetail @Summary 获取按钮详情
// @Description 获取指定按钮的详细信息
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param id path int true "按钮ID"
// @Success 200 {object} models.Button "按钮详情"
// @Failure 400 {object} ErrorResponse "无效的按钮ID"
// @Failure 500 {object} ErrorResponse "获取按钮详情失败"
// @Security ApiKeyAuth
// @Router /api/buttons/detail/{id} [get]
func (c *ButtonController) GetDetail(ctx *gin.Context) {
	buttonID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的按钮ID"})
		return
	}

	button, err := c.buttonService.GetButtonByID(buttonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取按钮详情失败"})
		return
	}

	ctx.JSON(http.StatusOK, button)
}

// GetButtonPermissionsResponse 获取按钮权限列表响应
type GetButtonPermissionsResponse []models.Permission

// CreateButton @Summary 创建按钮
// @Description 创建新按钮，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param button body CreateButtonRequest true "按钮信息"
// @Success 201 {object} CreateButtonResponse "按钮创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "创建按钮失败"
// @Security ApiKeyAuth
// @Router /api/buttons [post]
// CreateButton 创建按钮
func (c *ButtonController) CreateButton(ctx *gin.Context) {
	var button struct {
		Name   string `json:"name" binding:"required"`
		Code   string `json:"code" binding:"required"`
		MenuID int    `json:"menu_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&button); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	newButton := &models.Button{
		Name:           button.Name,
		PermissionCode: button.Code,
		MenuID:         button.MenuID,
	}

	if err := c.buttonService.CreateButton(newButton); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建按钮失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "按钮创建成功", "id": newButton.ID})
}

// UpdateButton @Summary 更新按钮信息
// @Description 更新指定按钮的信息，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param id path int true "按钮ID"
// @Param button body UpdateButtonRequest true "按钮信息"
// @Success 200 {object} UpdateButtonResponse "按钮信息更新成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "更新按钮信息失败"
// @Security ApiKeyAuth
// @Router /api/buttons/{id} [put]
// UpdateButton 更新按钮信息
func (c *ButtonController) UpdateButton(ctx *gin.Context) {
	buttonID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的按钮ID"})
		return
	}

	var updateData struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		MenuID int    `json:"menu_id"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	button := &models.Button{
		ID:             buttonID,
		Name:           updateData.Name,
		PermissionCode: updateData.Code,
		MenuID:         updateData.MenuID,
	}

	if err := c.buttonService.UpdateButton(button); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新按钮信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "按钮信息更新成功"})
}

// GetButtonsByMenuID @Summary 获取菜单按钮列表
// @Description 获取指定菜单的所有按钮列表，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param menuId path int true "菜单ID"
// @Success 200 {array} GetButtonsByMenuIDResponse "按钮列表"
// @Failure 400 {object} ErrorResponse "无效的菜单ID"
// @Failure 500 {object} ErrorResponse "获取按钮列表失败"
// @Security ApiKeyAuth
// @Router /api/buttons/menu/{menuId} [get]
// GetButtonsByMenuID 获取指定菜单的按钮列表
func (c *ButtonController) GetButtonsByMenuID(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("menuId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜单ID"})
		return
	}

	buttons, err := c.buttonService.GetButtonsByMenuID(menuID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取按钮列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, buttons)
}

// GetButtonsRequest 获取按钮列表请求参数
type GetButtonsRequest struct {
	Page     int `json:"page" example:"1" binding:"required"`
	PageSize int `json:"pageSize" example:"10" binding:"required"`
}

// GetButtonsResponse 按钮列表响应
type GetButtonsResponse struct {
	Data     []models.Button `json:"data"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// ListButtons @Summary 获取按钮列表
// @Description 获取所有按钮的列表，支持分页，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param request body GetButtonsRequest true "分页参数"
// @Success 200 {object} GetButtonsResponse "按钮列表"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "获取按钮列表失败"
// @Security ApiKeyAuth
// @Router /api/buttons [post]
// ListButtons 获取按钮列表
func (c *ButtonController) ListButtons(ctx *gin.Context) {
	var req GetButtonsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	buttons, total, err := c.buttonService.ListButtons(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取按钮列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     buttons,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// BindPermission @Summary 为按钮绑定权限
// @Description 为指定按钮绑定权限，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param id path int true "按钮ID"
// @Param permission body BindButtonPermissionRequest true "权限信息"
// @Success 200 {object} BindButtonPermissionResponse "权限绑定成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "绑定权限失败"
// @Security ApiKeyAuth
// @Router /api/buttons/{id}/permission [post]
// BindPermission 为按钮绑定权限
func (c *ButtonController) BindPermission(ctx *gin.Context) {
	buttonID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的按钮ID"})
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

	if err := c.buttonService.BindButtonPermission(buttonID, permission.PermissionCode, permission.Name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "绑定权限失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "权限绑定成功"})
}

// GetButtonPermissions @Summary 获取按钮权限列表
// @Description 获取指定按钮的权限列表，需要管理员权限
// @Tags 按钮管理
// @Accept json
// @Produce json
// @Param id path int true "按钮ID"
// @Success 200 {array} GetButtonPermissionsResponse "权限列表"
// @Failure 400 {object} ErrorResponse "无效的按钮ID"
// @Failure 500 {object} ErrorResponse "获取按钮权限列表失败"
// @Security ApiKeyAuth
// @Router /api/buttons/{id}/permissions [get]
// GetButtonPermissions 获取按钮的权限列表
func (c *ButtonController) GetButtonPermissions(ctx *gin.Context) {
	buttonID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的按钮ID"})
		return
	}

	permissions, err := c.buttonService.GetButtonPermissions(buttonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取按钮权限列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, permissions)
}
