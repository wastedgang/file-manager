package myspace

import (
	"fmt"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (m *MySpaceController) Download() gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		var err error
		currentUser := m.UserService.GetCurrentUser(ctx)
		// 检查文件是否存在
		downloadFilePath := "/" + ctx.Param("download_file_path")
		downloadFilePath, err = filepath.Abs(downloadFilePath)
		if err != nil {
			ctx.String(404, "文件不存在")
			return
		}
		fileInfo := m.MySpaceService.Get(currentUser.Id, downloadFilePath)
		if fileInfo == nil {
			ctx.String(404, "文件不存在")
			return
		}

		// 文件直接返回，文件夹要先压缩再返回
		if fileInfo.Type == fileinfotype.Normal {
			// 获取真正文件位置
			storeFileInfo := m.StoreFileService.Get(fileInfo.ContentHash)
			if storeFileInfo == nil {
				ctx.String(404, "文件不存在")
				return
			}
			storeFilePath := filepath.Join(storeFileInfo.StoreDirectoryPath, storeFileInfo.StoreFilename)
			file, err := os.Open(storeFilePath)
			if err != nil {
				ctx.String(500, "文件打开失败")
				return
			}
			defer file.Close()

			var startRange, endRange, contentLength int64 = 0, fileInfo.FileSize - 1, fileInfo.FileSize
			if rangeHeader := ctx.GetHeader("Range"); rangeHeader != "" {
				// 检查Range头部格式
				ranges := strings.Split(rangeHeader[6:], "-")
				if !strings.HasPrefix(rangeHeader, "bytes=") || len(ranges) == 0 || len(ranges) > 2 {
					ctx.String(400, "文件内容范围错误")
					return
				}

				// 计算Range开始范围
				if ranges[0] != "" {
					startRange, err = strconv.ParseInt(ranges[0], 10, 64)
					if err != nil {
						ctx.String(400, "文件内容范围错误")
						return
					}
				}
				// 计算Range结束范围
				if len(ranges) > 1 && ranges[1] != "" {
					endRange, err = strconv.ParseInt(ranges[1], 10, 64)
					if err != nil {
						ctx.String(400, "文件内容范围错误")
						return
					}
				}
				// 检查范围开始与结束的合法性
				if startRange > endRange {
					ctx.String(400, "文件内容范围错误")
					return
				}

				contentLength = endRange - startRange + 1

				// 设置文件开始读取位置
				_, err = file.Seek(startRange, io.SeekStart)
				if err != nil {
					ctx.String(500, "文件读取错误: %v", err)
					return
				}

				// 设置状态码、Content-Range与Accept-Ranges
				ctx.Status(206)
				ctx.Header("Accept-Ranges", "bytes")
				ctx.Header("Content-Range",
					fmt.Sprintf("bytes %d-%d/%d", startRange, endRange, fileInfo.FileSize))
			}

			// 设置Content-Type与Content-Length
			ctx.Header("Content-Type", fileInfo.MimeType)
			ctx.Header("Content-Length", strconv.FormatInt(contentLength, 10))

			// 传输文件
			io.CopyN(ctx.Writer, file, contentLength)
		} else {
			ctx.String(400, "暂不支持文件夹下载或多文件下载")
			return
			// TODO: 文件夹压缩返回

			// TODO: 压缩文件忽略断点续传
		}
	}
	return handler
}
