package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tenant-center/models"
	"tenant-center/services"
)

// @title 用户管理API
// @version 1.0
// @description 用户管理相关的API接口，包括用户登录、创建用户、更新用户信息和角色绑定等功能

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

// Login @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码，返回JWT token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param loginData body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse "登录成功，返回token"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 401 {object} ErrorResponse "用户名或密码错误"
// @Router /api/login [post]

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`  // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT令牌
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error" example:"无效的请求参数"` // 错误信息
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	token, err := c.userService.Login(loginData.Username, loginData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// CreateUser @Summary 创建用户
// @Description 创建新用户，需要管理员权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 201 {object} CreateUserResponse "用户创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "创建用户失败"
// @Security ApiKeyAuth
// @Router /api/users [post]

// UpdateUser @Summary 更新用户信息
// @Description 更新指定用户的信息，需要管理员权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body UpdateUserRequest true "用户信息"
// @Success 200 {object} UpdateUserResponse "用户信息更新成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "更新用户信息失败"
// @Security ApiKeyAuth
// @Router /api/users/{id} [put]

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	newUser := &models.User{
		Username: user.Username,
		Password: user.Password,
	}

	if err := c.userService.CreateUser(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "用户创建成功", "id": newUser.ID})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var updateData struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	user := &models.User{
		ID:       userID,
		Username: updateData.Username,
		Password: updateData.Password,
	}

	if err := c.userService.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// BindRolesRequest 绑定角色请求参数
type BindRolesRequest []int

// BindRolesResponse 绑定角色响应
type BindRolesResponse struct {
	Message string `json:"message" example:"角色绑定成功"` // 响应消息
}

// RouteResponse 路由响应数据
type RouteResponse []services.RouteItem

// BindRoles @Summary 为用户绑定角色
// @Description 为指定用户绑定一个或多个角色，需要管理员权限
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param roleIDs body BindRolesRequest true "角色ID列表" example:[1,2,3]
// @Success 200 {object} BindRolesResponse "角色绑定成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "绑定角色失败"
// @Security ApiKeyAuth
// @Router /api/users/{id}/roles [post]
// BindRoles 为用户绑定角色
func (c *UserController) BindRoles(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var roleIDs []int
	if err := ctx.ShouldBindJSON(&roleIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := c.userService.BindUserRoles(userID, roleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "绑定角色失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "角色绑定成功"})
}

// PageUsers @Summary 获取用户列表
// @Description 获取用户列表，支持分页
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body GetUsersRequest true "分页参数"
// @Success 200 {object} GetUsersResponse "用户列表"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "获取用户列表失败"
// @Security ApiKeyAuth
// @Router /api/users [post]
type GetUsersRequest struct {
	Page     int `json:"page" example:"1" binding:"required"`
	PageSize int `json:"pageSize" example:"10" binding:"required"`
}

func (c *UserController) PageUsers(ctx *gin.Context) {
	var req GetUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("绑定请求参数失败: %v\n", err) // 添加错误日志
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	users, total, err := c.userService.PageUsers(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     users,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// GetRoutes @Summary 获取用户路由数据
// @Description 获取用户的路由数据，包括菜单和权限信息
// @Tags 用户管理
// @Produce json
// @Success 200 {object} RouteResponse "路由数据"
// @Failure 400 {object} ErrorResponse "无效的用户ID"
// @Failure 500 {object} ErrorResponse "获取路由数据失败"
// @Security ApiKeyAuth
// @Router /api/users/routes [get]
func (c *UserController) GetRoutes(ctx *gin.Context) {
	// 从JWT中获取用户ID
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	routes, err := c.userService.GetUserRoutes(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取路由数据失败"})
		return
	}

	ctx.JSON(http.StatusOK, routes)
}
