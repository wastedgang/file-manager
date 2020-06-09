package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/model/constant/fileactiontype"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

func (m *MySpaceController) Copy() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		currentUser := m.UserService.GetCurrentUser(ctx)
		var err error
		var form struct {
			SourceDirectoryPath string `form:"source_directory_path" binding:"required,min=1"`
			Filename            string `form:"filename" binding:"required,min=2"`
			TargetDirectoryPath string `form:"target_directory_path" binding:"required,min=1"`
			ActionType          int    `form:"action" binding:"required,oneof=rename override"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}
		if form.SourceDirectoryPath == "" || strings.HasSuffix(form.SourceDirectoryPath, "/") && form.SourceDirectoryPath != "/" {
			return BadRequest
		}
		if form.TargetDirectoryPath == "" || strings.HasSuffix(form.TargetDirectoryPath, "/") && form.TargetDirectoryPath != "/" {
			return BadRequest
		}

		actionType := fileactiontype.FileActionType(form.ActionType)
		// 检查目录是否存在
		if !m.FileInfoService.IsDirectoryExists(currentUser.Id, form.SourceDirectoryPath) ||
			!m.FileInfoService.IsDirectoryExists(currentUser.Id, form.TargetDirectoryPath) {
			return DirectoryNotExists
		}

		// 获取欲复制的文件信息
		sourceFileInfo := m.FileInfoService.Get(currentUser.Id, filepath.Join(form.SourceDirectoryPath, form.Filename))
		if sourceFileInfo == nil {
			return FileNotExists
		}
		// 检查目标文件夹是否在欲移动文件夹里
		if sourceFileInfo.Type == fileinfotype.Directory {
			directoryPath := filepath.Join(sourceFileInfo.DirectoryPath, sourceFileInfo.Filename)
			if strings.HasPrefix(form.TargetDirectoryPath, directoryPath) {
				return TargetFolderInsideSourceFolder
			}
		}

		// 复制
		err = m.FileInfoService.Copy(sourceFileInfo, form.TargetDirectoryPath, actionType)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}
