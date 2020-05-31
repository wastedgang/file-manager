package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

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
		fileInfo := m.FileInfoService.Get(currentUser.Id, filepath.Join(directoryPath, form.Filename))
		if fileInfo != nil {
			if fileInfo.Type == fileinfotype.Directory {
				return Success
			}
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

func (m *MySpaceController) ListDirectories() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		currentUser := m.UserService.GetCurrentUser(ctx)
		fileInfos := m.FileInfoService.ListDirectories(currentUser.Id)

		// 构造目录信息表
		directoryMap := make(map[string]*model.DirectoryInfo)
		// 添加根目录信息
		directoryMap["/"] = &model.DirectoryInfo{
			Filename:            "/",
			Filepath:            "/",
			ParentDirectoryPath: "",
			HasSubDirectories:   false,
		}
		for _, fileInfo := range fileInfos {
			directoryInfo := &model.DirectoryInfo{
				Filename:            fileInfo.Filename,
				Filepath:            filepath.Join(fileInfo.DirectoryPath, fileInfo.Filename),
				ParentDirectoryPath: fileInfo.DirectoryPath,
				HasSubDirectories:   false,
			}
			directoryMap[directoryInfo.Filepath] = directoryInfo
		}
		// 重新计算HasSubDirectories字段值
		for _, fileInfo := range fileInfos {
			if directoryMap[fileInfo.DirectoryPath] == nil {
				continue
			}
			directoryMap[fileInfo.DirectoryPath].HasSubDirectories = true
		}

		directories := make([]*model.DirectoryInfo, 0, len(directoryMap))
		for _, directoryInfo := range directoryMap {
			directories = append(directories, directoryInfo)
		}
		return Success.AddField("directories", directories)
	})
	return handler
}
