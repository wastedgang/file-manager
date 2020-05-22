package controller

import (
	"github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/model/constant/usertype"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func RequireLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.GetCurrentUser(ctx)
		if user == nil {
			ctx.JSON(200, statuscode.NotLogin)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func RequireSystemAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.GetCurrentUser(ctx)
		if user == nil {
			ctx.JSON(200, statuscode.NotLogin)
			ctx.Abort()
			return
		}
		if user.Type != usertype.SystemAdmin {
			ctx.JSON(200, statuscode.PermissionDenied)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
