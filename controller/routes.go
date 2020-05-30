package controller

import (
	"github.com/farseer810/file-manager/controller/auth"
	"github.com/farseer810/file-manager/controller/group"
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
	var groupController *group.GroupController
	err := inject.Get(
		&settingsController,
		&userController,
		&authController,
		&storeSpaceController,
		&mySpaceController,
		&groupController,
	)
	if err != nil {
		return err
	}

	// 配置中间件
	r.Use(middleware.RecoveryHandler())
	r.Use(middleware.CheckConfigurationHandler())
	r.Use(middleware.JwtHandler())

	v1 := r.Group("/api/v1")
	{
		// 配置信息管理
		v1.GET("/settings/database", RequireLogin(), settingsController.GetDatabaseSettings())
		v1.POST("/settings/database", settingsController.SetDatabaseSettings())
		v1.POST("/settings/database/check", settingsController.CheckDatabaseSettings())

		// 授权管理
		v1.POST("/auth/login", authController.Login())
		v1.Any("/auth/logout", RequireLogin(), authController.Logout())

		// 用户管理
		v1.POST("/user", RequireSystemAdmin(), userController.AddUser())
		v1.PUT("/user/:username/info", RequireSystemAdmin(), userController.UpdateUserInfo())
		v1.PATCH("/user/:username/password", RequireSystemAdmin(), userController.UpdateUserPassword())
		v1.DELETE("/user/:username", RequireSystemAdmin(), userController.DeleteUser())
		v1.GET("/users", RequireLogin(), userController.ListUsers())
		v1.GET("/user/:username/info", RequireSystemAdmin(), userController.GetUserInfo())
		v1.GET("/current_user/info", RequireLogin(), userController.GetCurrentUser())
		v1.PUT("/current_user/info", RequireLogin(), userController.UpdateCurrentUserInfo())
		v1.PATCH("/current_user/password", RequireLogin(), userController.UpdateCurrentUserPassword())

		// 存储空间管理
		v1.POST("/store_space", RequireSystemAdmin(), storeSpaceController.AddStoreSpace())
		v1.PUT("/store_space", RequireSystemAdmin(), storeSpaceController.UpdateStoreSpace())
		v1.DELETE("/store_space", RequireSystemAdmin(), storeSpaceController.DeleteStoreSpace())
		v1.GET("/store_spaces", RequireSystemAdmin(), storeSpaceController.ListStoreSpaces())
		v1.GET("/store_space/files", RequireSystemAdmin(), storeSpaceController.ListStoreSpaceFiles())
		v1.GET("/directories", storeSpaceController.ListDirectories())

		// 个人空间管理
		v1.POST("/my_space/directory", RequireLogin(), mySpaceController.AddDirectory())
		v1.GET("/my_space/directories", RequireLogin(), mySpaceController.ListDirectories())
		v1.GET("/my_space/files", RequireLogin(), mySpaceController.List())
		v1.DELETE("/my_space/file", RequireLogin(), mySpaceController.Delete())
		v1.PUT("/my_space/file", RequireLogin(), mySpaceController.Rename())
		v1.POST("/my_space/file/copy", RequireLogin(), mySpaceController.Copy())
		v1.POST("/my_space/file/move", RequireLogin(), mySpaceController.Move())
		v1.POST("/my_space/file/share", RequireLogin(), mySpaceController.Share())

		// 个人空间上传
		v1.POST("/my_space/file/upload/*upload_directory_path", RequireLogin(), mySpaceController.Upload())

		// 个人空间下载
		v1.GET("/my_space/file/download/*download_file_path", RequireLogin(), mySpaceController.Download())

		// 群组管理
		v1.POST("/group", RequireLogin(), groupController.Add())
		v1.PUT("/group/:group_name", RequireLogin(), groupController.Update())
		//v1.GET("/group/:group_name", RequireLogin(), groupController.Get())
		v1.DELETE("/group/:group_name", RequireLogin(), groupController.Delete())
		v1.GET("/user_groups", RequireLogin(), groupController.ListUserGroup())
		//v1.GET("/user_group/:group_name", RequireLogin(), groupController.GetUserGroup())

		// 群组成员管理
		v1.POST("/group/:group_name/member", RequireLogin(), groupController.AddMember())
		v1.PUT("/group/:group_name/member", RequireLogin(), groupController.UpdateMember())
		v1.DELETE("/group/:group_name/member/:user_id", RequireLogin(), groupController.DeleteMember())
		v1.GET("/group/:group_name/members", RequireLogin(), groupController.ListMembers())
	}

	return nil
}
