package statuscode

import (
	"github.com/farseer810/file-manager/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	utils.RegisterTimeSerializer(utils.DefaultTimeFormat, utils.DefaultTimeLocation)
}

type HandlerFunc func(ctx *gin.Context) *Response

func ConvertGinHandlerFunc(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				response, ok := r.(*Response)
				if !ok {
					panic(r)
				}
				ctx.Render(200, JSONRender{response})
			}

		}()
		response := handler(ctx)
		ctx.Render(200, JSONRender{response})
	}
}
