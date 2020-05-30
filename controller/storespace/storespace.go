package storespace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/service"
	"github.com/farseer810/file-manager/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func init() {
	inject.Provide(new(StoreSpaceController))
}

type StoreSpaceController struct {
	StoreSpaceService *service.StoreSpaceService
	StoreFileService  *service.StoreFileService
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

// UpdateStoreSpace 更新存储空间
func (s *StoreSpaceController) UpdateStoreSpace() gin.HandlerFunc {
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

		// 检查存储空间是否存在
		storeSpace := s.StoreSpaceService.GetByDirectoryPath(form.DirectoryPath)
		if storeSpace == nil {
			return StoreSpaceNotExists
		}

		err = s.StoreSpaceService.Update(storeSpace.DirectoryPath, form.AllocateSize, form.Remark)
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

// ListStoreSpaceFiles 获取存储空间文件列表
func (s *StoreSpaceController) ListStoreSpaceFiles() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		directoryPath := ctx.DefaultQuery("directory_path", "")
		if directoryPath == "" {
			return BadRequest
		}

		// 检查存储空间是否存在
		storeSpace := s.StoreSpaceService.GetByDirectoryPath(directoryPath)
		if storeSpace == nil {
			return StoreSpaceNotExists
		}

		storeFiles := s.StoreFileService.ListByStoreDirectoryPath(storeSpace.DirectoryPath)

		return Success.AddField("store_files", storeFiles)
	})
	return handler
}

func (s *StoreSpaceController) ListDirectories() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		directoryPath := ctx.DefaultQuery("directory_path", "")
		if directoryPath == "" || !strings.HasPrefix(directoryPath, "/") {
			return BadRequest
		}
		directoryPath = filepath.Clean(directoryPath)

		// 检查目录是否存在
		if !utils.IsDir(directoryPath) {
			return DirectoryNotExists
		}

		fileInfos, err := ioutil.ReadDir(directoryPath)
		if err != nil {
			return InternalServerError
		}

		directories := make([]*model.DirectoryInfo, 0)
		for _, fileInfo := range fileInfos {
			filename := fileInfo.Name()
			if !fileInfo.IsDir() || strings.HasPrefix(filename, ".") {
				continue
			}

			// 是否有子文件夹
			hasSubDirectories := false
			subFileInfos, _ := ioutil.ReadDir(filepath.Join(directoryPath, filename))
			for _, subFileInfo := range subFileInfos {
				if subFileInfo.IsDir() && !strings.HasPrefix(subFileInfo.Name(), ".") {
					hasSubDirectories = true
				}
			}

			directoryInfo := &model.DirectoryInfo{
				Filename:          filename,
				Filepath:          filepath.Join(directoryPath, filename),
				ParentDirectoryPath: directoryPath,
				HasSubDirectories: hasSubDirectories,
			}
			directories = append(directories, directoryInfo)
		}
		return Success.AddField("directories", directories)
	})
	return handler
}
