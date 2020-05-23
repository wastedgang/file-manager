package controller

import (
	"github.com/farseer810/file-manager/controller/middleware"
	"github.com/farseer810/file-manager/inject"
	"github.com/gin-gonic/gin"
)

// InitRoutes 配置路由
func InitRoutes(r *gin.Engine) error {
	var settingsController *SettingsController
	var userController *UserController
	var authController *AuthController
	var storeSpaceController *StoreSpaceController
	err := inject.Get(
		&settingsController,
		&userController,
		&authController,
		&storeSpaceController,
	)
	if err != nil {
		return err
	}

	// 配置中间件
	r.Use(middleware.RecoveryHandler())
	r.Use(middleware.JwtHandler())
	r.Use(middleware.CheckConfigurationHandler())

	v1 := r.Group("/api/v1")
	{
		// 配置信息管理
		v1.GET("/settings/database", RequireSystemAdmin(), settingsController.GetDatabaseSettings())
		v1.POST("/settings/database", settingsController.SetDatabaseSettings())
		v1.POST("/settings/database/check", settingsController.CheckDatabaseSettings())

		// 授权管理
		v1.POST("/auth/login", authController.Login())
		v1.Any("/auth/logout", RequireLogin(), authController.Logout())

		// 用户管理
		v1.POST("/user", RequireSystemAdmin(), userController.AddUser())
		v1.PUT("/user/:user_id/info", RequireSystemAdmin(), userController.UpdateUserInfo())
		v1.PATCH("/user/:user_id/password", RequireSystemAdmin(), userController.UpdateUserPassword())
		v1.DELETE("/user/:user_id", RequireSystemAdmin(), userController.DeleteUser())
		v1.GET("/users", RequireLogin(), userController.ListUsers())
		v1.GET("/current_user/info", RequireSystemAdmin(), userController.GetCurrentUser())
		v1.PUT("/current_user/info", RequireSystemAdmin(), userController.UpdateCurrentUserInfo())
		v1.PATCH("/current_user/password", RequireSystemAdmin(), userController.UpdateCurrentUserPassword())

		// 存储空间管理
		v1.POST("/store_space", RequireSystemAdmin(), storeSpaceController.AddStoreSpace())
		v1.DELETE("/store_space", RequireSystemAdmin(), storeSpaceController.DeleteStoreSpace())
		v1.POST("/store_spaces", RequireSystemAdmin(), storeSpaceController.ListStoreSpaces())
	}

	return nil
}
