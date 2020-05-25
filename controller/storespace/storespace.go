package storespace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/service"
	"github.com/farseer810/file-manager/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

func init() {
	inject.Provide(new(StoreSpaceController))
}

type StoreSpaceController struct {
	StoreSpaceService *service.StoreSpaceService
}

// AddStoreSpace 添加存储空间
func (s *StoreSpaceController) AddStoreSpace() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			DirectoryPath string `form:"directory_path" binding:"required,min=1"`
			AllocateSize  int64  `form:"allocate_size" binding:"required,min=10485760"`
			Remark        string `form:"remark" binding:"max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}
		// 保证是绝对路径
		if !strings.HasPrefix(form.DirectoryPath, "/") {
			return BadRequest
		}

		// 目录是否存在
		directoryPath := filepath.Clean(form.DirectoryPath)
		if !utils.IsDir(directoryPath) {
			return DirectoryNotExists
		}

		// 检查存储空间是否存在
		storeSpace := s.StoreSpaceService.GetByDirectoryPath(form.DirectoryPath)
		if storeSpace != nil {
			return StoreSpaceExists
		}

		err = s.StoreSpaceService.Add(directoryPath, form.AllocateSize, form.Remark)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// DeleteStoreSpace 删除存储空间
func (s *StoreSpaceController) DeleteStoreSpace() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		directoryPath := ctx.DefaultQuery("directory_path", "")
		if directoryPath == "" {
			return BadRequest
		}

		// 检查存储空间是否存在
		storeSpace := s.StoreSpaceService.GetByDirectoryPath(directoryPath)
		if storeSpace == nil {
			return StoreSpaceNotExists
		}

		err = s.StoreSpaceService.Delete(directoryPath)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

// ListStoreSpaces 获取存储空间列表
func (s *StoreSpaceController) ListStoreSpaces() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		storeSpaces := s.StoreSpaceService.List()
		return Success.AddField("store_spaces", storeSpaces)
	})
	return handler
}
