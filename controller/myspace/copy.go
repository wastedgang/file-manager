package myspace

import (
	"encoding/json"
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/model"
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
			Filenames           string `form:"filenames" binding:"required,min=2"`
			TargetDirectoryPath string `form:"target_directory_path" binding:"required,min=1"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}
		var filenames []string
		if err = json.Unmarshal([]byte(form.Filenames), &filenames); err != nil {
			return BadRequest
		}
		if len(filenames) == 0 {
			return Success
		}

		// 检查目录是否存在
		if !m.FileInfoService.IsDirectoryExists(currentUser.Id, form.SourceDirectoryPath) ||
			!m.FileInfoService.IsDirectoryExists(currentUser.Id, form.TargetDirectoryPath) {
			return DirectoryNotExists
		}

		// 检查欲复制的文件列表
		fileInfos := m.FileInfoService.ListByFilenames(currentUser.Id, form.SourceDirectoryPath, filenames)
		if len(fileInfos) == 0 {
			return Success
		}
		// 检查目标文件夹是否在欲复制文件夹里
		var directoryInfos []*model.FileInfo
		for _, fileInfo := range fileInfos {
			if fileInfo.Type != fileinfotype.Directory {
				continue
			}
			directoryInfos = append(directoryInfos, fileInfo)
		}
		if len(directoryInfos) != 0 {
			for _, directoryInfo := range directoryInfos {
				directoryPath := filepath.Join(directoryInfo.DirectoryPath, directoryInfo.Filename)
				if strings.HasPrefix(form.TargetDirectoryPath, directoryPath) {
					return TargetFolderInsideSourceFolder
				}
			}
		}

		// 复制
		err = m.FileInfoService.Copy(currentUser.Id, form.SourceDirectoryPath, filenames, form.TargetDirectoryPath)
		if err != nil {
			return InternalServerError
		}
		return Success
	})
	return handler
}
