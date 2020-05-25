package controller

import (
	"github.com/farseer810/file-manager/controller/auth"
	"github.com/farseer810/file-manager/controller/middleware"
	"github.com/farseer810/file-manager/controller/myspace"
	"github.com/farseer810/file-manager/controller/settings"
	"github.com/farseer810/file-manager/controller/storespace"
	"github.com/farseer810/file-manager/controller/user"
	"github.com/farseer810/file-manager/inject"
	"github.com/gin-gonic/gin"
)

// InitRoutes 配置路由
func InitRoutes(r *gin.Engine) error {
	var settingsController *settings.SettingsController
	var userController *user.UserController
	var authController *auth.AuthController
	var storeSpaceController *storespace.StoreSpaceController
	var mySpaceController *myspace.MySpaceController
	err := inject.Get(
		&settingsController,
		&userController,
		&authController,
		&storeSpaceController,
		&mySpaceController,
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
		v1.GET("/user/:user_id/info", RequireSystemAdmin(), userController.GetUserInfo())
		v1.GET("/current_user/info", RequireLogin(), userController.GetCurrentUser())
		v1.PUT("/current_user/info", RequireLogin(), userController.UpdateCurrentUserInfo())
		v1.PATCH("/current_user/password", RequireLogin(), userController.UpdateCurrentUserPassword())

		// 存储空间管理
		v1.POST("/store_space", RequireSystemAdmin(), storeSpaceController.AddStoreSpace())
		v1.DELETE("/store_space", RequireSystemAdmin(), storeSpaceController.DeleteStoreSpace())
		v1.GET("/store_spaces", RequireSystemAdmin(), storeSpaceController.ListStoreSpaces())

		// 个人空间管理
		v1.GET("/my_space/files", RequireLogin(), mySpaceController.List())
		v1.POST("/my_space/directory", RequireLogin(), mySpaceController.AddDirectory())
		v1.DELETE("/my_space/file", RequireLogin(), mySpaceController.Delete())
		v1.PUT("/my_space/file", RequireLogin(), mySpaceController.Rename())
		v1.POST("/my_space/file/copy", RequireLogin(), mySpaceController.Copy())
		v1.POST("/my_space/file/move", RequireLogin(), mySpaceController.Move())
		v1.POST("/my_space/file/share", RequireLogin(), mySpaceController.Share())

		// 个人空间上传
		v1.POST("/my_space/file/upload/*upload_directory_path", RequireLogin(), mySpaceController.Upload())
		v1.GET("/my_space/file/upload/start_point", RequireLogin(), mySpaceController.GetUploadStartPoint())

		// 个人空间下载
		v1.GET("/my_space/file/download/*download_file_path", RequireLogin(), mySpaceController.Download())
	}

	return nil
}
