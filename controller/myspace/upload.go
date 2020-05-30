package myspace

import (
	"fmt"
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"path/filepath"
)


func (m *MySpaceController) Upload() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		currentUser := m.UserService.GetCurrentUser(ctx)

		// 检查上传目录路径参数
		directoryPath := "/" + ctx.Param("upload_directory_path")
		directoryPath, err = filepath.Abs(directoryPath)
		if err != nil {
			return BadRequest
		}
		// 目录是否存在
		if !m.FileInfoService.IsDirectoryExists(currentUser.Id, directoryPath) {
			return DirectoryNotExists
		}

		// 检查是否存在存储空间
		fmt.Println(m.StoreSpaceService.List())
		if len(m.StoreSpaceService.List()) == 0 {
			return StoreSpaceNotExists.SetMessage("请先添加存储空间")
		}

		multipartReader, err := ctx.Request.MultipartReader()
		if err != nil {
			return BadRequest
		}
		var filePart *multipart.Part
		for {
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return InternalServerError
			}

			if part.FormName() == "file" {
				filePart = part
				break
			}
		}
		// 检查参数
		if filePart == nil {
			return BadRequest
		}

		err = m.StoreFileService.Save(currentUser.Id, filePart, directoryPath)
		if err != nil {
			return UploadFailed
		}
		return Success
	})
	return handler
}
