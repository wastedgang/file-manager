package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
)

func (m *MySpaceController) Share() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}
