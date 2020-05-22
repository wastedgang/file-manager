package middleware

import (
	"github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RecoveryHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				logrus.Error(r)
				ctx.JSON(200, statuscode.InternalServerError)
			}
		}()
		ctx.Next()
	}
}
