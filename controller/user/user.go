package user

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/usertype"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
	"regexp"
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
		var err error
		var form struct {
			Username string `form:"username" binding:"required"`
			Remark   string `form:"remark" binding:"max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查username参数
		if ok, _ := regexp.MatchString("^[a-zA-Z0-9_.@-]{6,32}$", form.Username); !ok {
			return BadRequest.SetMessage("用户名格式错误")
		}

		// 检查用户名是否已存在
		user := u.UserService.GetByUsername(form.Username)
		if user != nil {
			return UsernameExists
		}

		// 添加用户
		_, err = u.UserService.Add(form.Username, form.Remark)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// UpdateUserInfo 更新用户信息
func (u *UserController) UpdateUserInfo() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Nickname string `form:"nickname" binding:"required"`
			Remark   string `form:"remark" binding:"max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查用户是否存在
		username := ctx.Param("username")
		user := u.UserService.GetByUsername(username)
		if user == nil {
			return UserNotExists
		}

		// 检查权限
		currentUser := u.UserService.GetCurrentUser(ctx)
		if currentUser.Id != user.Id && user.Type == usertype.SystemAdmin {
			return PermissionDenied
		}

		err = u.UserService.Update(user.Id, form.Nickname, form.Remark)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// UpdateCurrentUserInfo 更新当前用户信息
func (u *UserController) UpdateCurrentUserInfo() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Nickname string `form:"nickname" binding:"required"`
			Remark   string `form:"remark" binding:"max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		currentUser := u.UserService.GetCurrentUser(ctx)
		err = u.UserService.Update(currentUser.Id, form.Nickname, form.Remark)
		if err != nil {
			return InternalServerError
		}

		currentUser = u.UserService.GetById(currentUser.Id)
		return Success.AddField("user", currentUser)
	})
	return handler
}

// UpdateUserPassword 更新用户密码
func (u *UserController) UpdateUserPassword() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Password string `form:"password" binding:"required,min=6"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查用户是否存在
		username := ctx.Param("username")
		user := u.UserService.GetByUsername(username)
		if user == nil {
			return UserNotExists
		}

		// 检查权限
		currentUser := u.UserService.GetCurrentUser(ctx)
		if currentUser.Id != user.Id && user.Type == usertype.SystemAdmin {
			return PermissionDenied
		}

		err = u.UserService.UpdateUserPassword(currentUser.Id, form.Password)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// UpdateCurrentUserPassword 更新当前用户密码
func (u *UserController) UpdateCurrentUserPassword() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			OldPassword string `form:"old_password" binding:"required"`
			NewPassword string `form:"new_password" binding:"required,min=6"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 检查权限
		currentUser := u.UserService.GetCurrentUser(ctx)
		if currentUser.Password != u.UserService.CalculateHashPassword(form.OldPassword) {
			return InvalidPassword
		}

		err = u.UserService.UpdateUserPassword(currentUser.Id, form.NewPassword)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// DeleteUser 删除用户
func (u *UserController) DeleteUser() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		// 检查用户是否存在
		username := ctx.Param("username")
		user := u.UserService.GetByUsername(username)
		if user == nil {
			return UserNotExists
		}

		// 检查权限
		currentUser := u.UserService.GetCurrentUser(ctx)
		if currentUser.Id == user.Id {
			return CannotDeleteCurrentUser
		}
		if user.Type == usertype.SystemAdmin && currentUser.Username != service.DefaultSystemAdminUsername {
			return PermissionDenied
		}

		err = u.UserService.DeleteById(user.Id)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// ListUsers 获取用户列表
func (u *UserController) ListUsers() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			SearchWord string `form:"search_word"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		users := u.UserService.List(form.SearchWord)
		return Success.AddField("users", users)
	})
	return handler
}

// GetCurrentUser 获取当前登录用户
func (u *UserController) GetCurrentUser() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		user := u.UserService.GetCurrentUser(ctx)
		user = u.UserService.GetById(user.Id)
		if user == nil {
			return UserNotExists
		}
		return Success.AddField("user", user)
	})
	return handler
}

// GetUserInfo 获取指定用户
func (u *UserController) GetUserInfo() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		// 检查用户是否存在
		username := ctx.Param("username")
		user := u.UserService.GetByUsername(username)
		if user == nil {
			return UserNotExists
		}
		return Success.AddField("user", user)
	})
	return handler
}
