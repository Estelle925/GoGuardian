package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tenant-center/controllers"
	"tenant-center/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建控制器实例
	userController := controllers.NewUserController(db)
	roleController := controllers.NewRoleController(db)
	permissionController := controllers.NewPermissionController(db)
	menuController := controllers.NewMenuController(db)
	buttonController := controllers.NewButtonController(db)

	// 公开路由组
	public := r.Group("/api")
	{
		// 用户登录
		public.POST("/login", userController.Login)
	}

	// 需要认证的路由组
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		user := protected.Group("/users")
		{
			user.POST("", userController.CreateUser)
			user.PUT("/:id", userController.UpdateUser)
			user.POST("/:id/roles", userController.BindRoles)
			user.GET("/routes", userController.GetRoutes)
			user.POST("/page", userController.PageUsers)
		}

		// 角色相关路由
		role := protected.Group("/roles")
		{
			role.POST("", roleController.CreateRole)
			role.PUT("/:id", roleController.UpdateRole)
			role.GET("detail/:id/", roleController.GetDetail)
			role.POST("/page", roleController.PageRoles)
			role.POST("/:id/bindPermissions", roleController.BindPermissions)
			role.GET("/:id/permissions", roleController.GetRolePermissions)
		}

		// 权限相关路由
		permission := protected.Group("/permissions")
		{
			permission.POST("", permissionController.CreatePermission)
			permission.GET("detail/:id/", permissionController.GetPermissionDetail)
			permission.PUT("/:id", permissionController.UpdatePermission)
			permission.POST("page", permissionController.PagePermissions)
			permission.GET("/type/:type", permissionController.GetPermissionsByType)
		}

		// 菜单相关路由
		menu := protected.Group("/menus")
		{
			menu.POST("", menuController.CreateMenu)
			menu.PUT("/:id", menuController.UpdateMenu)
			menu.GET("detail/:id/", menuController.GetDetail)
			menu.POST("/page", menuController.ListMenus)
			menu.GET("/parent/:parentId", menuController.GetMenusByParentID)
			menu.POST("/:id/permission", menuController.BindPermission)
			menu.GET("/:id/permissions", menuController.GetMenuPermissions)
		}

		// 按钮相关路由
		button := protected.Group("/buttons")
		{
			button.POST("", buttonController.CreateButton)
			button.PUT("/:id", buttonController.UpdateButton)
			button.GET("detail/:id/", buttonController.GetDetail)
			button.GET("/menu/:menuId", buttonController.GetButtonsByMenuID)
			button.POST("/:id/permission", buttonController.BindPermission)
			button.GET("/:id/permissions", buttonController.GetButtonPermissions)
			button.POST("/page", buttonController.ListButtons)
		}
	}
}
