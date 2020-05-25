package myspace

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
)

func (m *MySpaceController) GetUploadStartPoint() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		contentHash := ctx.Query("contentHash")
		if contentHash == "" {
			return BadRequest
		}

		fileSize, err := strconv.ParseInt(ctx.Query("file_size"), 10, 64)
		if err != nil {
			return BadRequest
		}

		// 检查是否已有存储空间
		bestStoreSpace := m.StoreSpaceService.GetBestStoreSpace()
		if bestStoreSpace == nil {
			return NoStoreSpace
		}
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
		// 检查Content-Hash是否存在
		contentHash := ctx.GetHeader("Content-Hash")
		if contentHash == "" {
			return BadRequest
		}

		currentUser := m.UserService.GetCurrentUser(ctx)

		// 内容开始位置
		contentStartAtHeader := ctx.GetHeader("Content-Start-At")
		if contentStartAtHeader == "" {
			contentStartAtHeader = "0"
		}
		contentStartAt, err := strconv.ParseInt(contentStartAtHeader, 10, 64)
		if err != nil {
			return BadRequest
		}
		// 检查开始位置是否正确
		startPoint := m.OngoingUploadService.GetUploadStartPoint(currentUser.Id, contentHash)
		storeFileInfo := m.StoreFileService.Get(contentHash)
		if storeFileInfo == nil && startPoint != contentStartAt {
			return InvalidUploadStartPoint
		}

		// 检查上传目录路径参数
		directoryPath := "/" + ctx.Param("upload_directory_path")
		directoryPath, err = filepath.Abs(directoryPath)
		if err != nil {
			return BadRequest
		}
		// 目录是否存在
		if !m.MySpaceService.IsDirectoryExists(currentUser.Id, directoryPath) {
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

		err = m.StoreFileService.Save(currentUser.Id, contentHash, filePart, directoryPath)
		if err != nil {
			return UploadFailed
		}
		return Success
	})
	return handler
}
