package controller

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.Provide(new(AuthController))
}

type AuthController struct {
	UserService *service.UserService
}

func (a *AuthController) Login() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Username string `form:"username" binding:"required"`
			Password string `form:"password" binding:"max=64"`
			Source   string `form:"source" binding:"max=64"`
		}
		if err := ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 获取用户信息
		user := a.UserService.GetByUsername(form.Username)

		// 用户不存在
		if user == nil {
			if form.Username == service.DefaultSystemAdminUsername {
				// 添加默认系统管理员
				user, err = a.UserService.AddDefaultSystemAdmin()
				if err != nil {
					return InternalServerError
				}
			} else {
				return UserNotExists
			}
		}

		// 计算哈希密码
		hashPassword := a.UserService.CalculateHashPassword(form.Password)
		// 密码错误
		if hashPassword != user.Password {
			return InvalidUsernameOrPassword
		}

		// 添加登录记录
		err = a.UserService.AddLoginRecord(user, form.Source)
		if err != nil {
			return InternalServerError
		}

		// 保存jwt到cookies
		token, err := a.UserService.GenerateToken(user)
		if err != nil {
			return InternalServerError
		}

		ctx.SetCookie(service.TokenCookieKey, token, 0, "/", "", false, true)
		return Success
	})
	return handler
}

func (a *AuthController) Logout() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		ctx.SetCookie(service.TokenCookieKey, "", -1, "/", "", false, true)
		return Success
	})
	return handler
}
