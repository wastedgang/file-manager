package controller

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.Provide(new(UserController))
}

type UserController struct {
	UserService *service.UserService
}

// AddUser 添加用户
func (u *UserController) AddUser() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// UpdateUserInfo 更新用户信息
func (u *UserController) UpdateUserInfo() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// UpdateUserPassword 更新用户密码
func (u *UserController) UpdateUserPassword() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// DeleteUser 删除用户
func (u *UserController) DeleteUser() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// ListUsers 获取用户列表
func (u *UserController) ListUsers() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// GetCurrentUser 获取当前登录用户
func (u *UserController) GetCurrentUser() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		user := u.UserService.GetCurrentUser(ctx)
		return Success.AddField("user", user)
	})
	return handler
}
