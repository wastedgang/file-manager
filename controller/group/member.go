package group

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
)

func (g *GroupController) AddMember() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

func (g *GroupController) UpdateMember() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

func (g *GroupController) DeleteMember() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

func (g *GroupController) ListMembers() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}
