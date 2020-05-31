package myspace

import (
	"fmt"
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/farseer810/file-manager/service"
	"github.com/farseer810/file-manager/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func init() {
	inject.Provide(new(MySpaceController))
}

type MySpaceController struct {
	UserService          *service.UserService
	FileInfoService      *service.FileInfoService
	StoreSpaceService    *service.StoreSpaceService
	StoreFileService     *service.StoreFileService
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
		//sort.SliceStable(fileInfos, func(i, j int) bool {
		//	return fileInfos[i].Less(fileInfos[j])
		//})
		for _, fileInfo := range fileInfos {
			fmt.Println(fileInfo.Filename, utils.GetFirstPinYinLetter(fileInfo.Filename))
		}
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
