package middleware

import (
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func JwtHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.ParseUserFromJwt(ctx)
		if user != nil {
			ctx.Set(service.CurrentUserContextName, user)
		}
		ctx.Next()
	}
}
