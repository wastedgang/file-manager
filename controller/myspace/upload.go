package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
)

// GetUploadStartPoint 获取上传开始位置
func (m *MySpaceController) GetUploadStartPoint() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		contentHash := ctx.Query("content_hash")
		if contentHash == "" {
			return BadRequest
		}

		fileSize, err := strconv.ParseInt(ctx.Query("file_size"), 10, 64)
		if err != nil {
			return BadRequest
		}

		storeFileInfo := m.StoreFileService.Get(contentHash)
		if storeFileInfo != nil {
			return Success.AddField("upload_start_point", storeFileInfo.FileSize)
		}

		// 检查是否已有存储空间
		bestStoreSpace := m.StoreSpaceService.GetBestStoreSpace()
		if bestStoreSpace == nil {
			return NoStoreSpace
		}
		// 检查剩余空间是否足够
		if bestStoreSpace.TotalFreeSpace < fileSize {
			return NotEnoughFreeSpace
		}

		currentUser := m.UserService.GetCurrentUser(ctx)
		startPoint := m.OngoingUploadService.GetUploadStartPoint(currentUser.Id, contentHash)
		return Success.AddField("upload_start_point", startPoint)
	})
	return handler
}

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
