package controller

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.Provide(new(StoreSpaceController))
}

type StoreSpaceController struct{}

// AddStoreSpace 添加存储空间
func (s *StoreSpaceController) AddStoreSpace() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// DeleteStoreSpace 删除存储空间
func (s *StoreSpaceController) DeleteStoreSpace() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

// ListStoreSpaces 获取存储空间列表
func (s *StoreSpaceController) ListStoreSpaces() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}
