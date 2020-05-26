package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

func init() {
	inject.Provide(new(MySpaceController))
}

type MySpaceController struct {
	UserService          *service.UserService
	FileInfoService      *service.FileInfoService
	StoreSpaceService    *service.StoreSpaceService
	StoreFileService     *service.StoreFileService
	OngoingUploadService *service.OngoingUploadService
}

func (m *MySpaceController) AddDirectory() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			DirectoryPath string `form:"directory_path"`
			Filename      string `form:"filename" binding:"required,min=1,max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}
		// 检查目录
		directoryPath := form.DirectoryPath
		if directoryPath == "" {
			directoryPath = "/"
		}
		if !strings.HasPrefix(directoryPath, "/") {
			return BadRequest
		}
		directoryPath, err = filepath.Abs(directoryPath)
		if err != nil {
			return BadRequest
		}

		// 检查文件夹名
		if form.Filename == "." || form.Filename == ".." || strings.Contains(form.Filename, "/") {
			return BadRequest
		}

		currentUser := m.UserService.GetCurrentUser(ctx)
		// 检查文件夹是否存在
		if !m.FileInfoService.IsDirectoryExists(currentUser.Id, directoryPath) {
			return DirectoryNotExists
		}

		// 检查文件是否已存在
		if m.FileInfoService.Get(currentUser.Id, filepath.Join(directoryPath, form.Filename)) != nil {
			return FileExists
		}

		err = m.FileInfoService.CreateDirectory(currentUser.Id, filepath.Join(directoryPath, form.Filename))
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}

func (m *MySpaceController) List() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		searchWord := ctx.DefaultQuery("search_word", "")
		directoryPath := ctx.DefaultQuery("directory_path", "/")
		if directoryPath == "" {
			directoryPath = "/"
		}
		currentUser := m.UserService.GetCurrentUser(ctx)
		fileInfos := m.FileInfoService.List(currentUser.Id, directoryPath, searchWord)
		return Success.AddField("files", fileInfos)
	})
	return handler
}

func (m *MySpaceController) Delete() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		return Success
	})
	return handler
}

func (m *MySpaceController) Rename() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		currentUser := m.UserService.GetCurrentUser(ctx)
		var err error
		var form struct {
			DirectoryPath string `form:"required,min=1,directory_path"`
			OldFilename   string `form:"old_filename" binding:"required,min=1,max=128"`
			NewFilename   string `form:"new_filename" binding:"required,min=1,max=128"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		// 文件是否存在
		fileInfo := m.FileInfoService.Get(currentUser.Id, filepath.Join(form.DirectoryPath, form.OldFilename))
		if fileInfo == nil {
			return FileNotExists
		}

		// 新文件名是否被占用
		if m.FileInfoService.Get(currentUser.Id, filepath.Join(form.DirectoryPath, form.NewFilename)) == nil {
			if fileInfo.Type == fileinfotype.Normal {
				return FileExists
			} else {
				return DirectoryExists
			}
		}

		err = m.FileInfoService.Rename(fileInfo, form.NewFilename)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}
