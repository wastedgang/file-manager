package middleware

import (
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.AddCallback(func(u *service.UserService) {
		userService = u
	})
}

var (
	userService *service.UserService
)

func JwtHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.ParseUserFromJwt(ctx)
		if user != nil {
			user = userService.GetById(user.Id)
			ctx.Set(service.CurrentUserContextName, user)
		}
		ctx.Next()
	}
}
