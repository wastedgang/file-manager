package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/farseer810/file-manager/utils"
	"github.com/jinzhu/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func init() {
	inject.Provide(new(StoreFileService))
}

type StoreFileService struct {
	OngoingUploadService *OngoingUploadService
	MySpaceService       *FileInfoService
	FileInfoService      *FileInfoService
}

func (s *StoreFileService) List() []*model.StoreFileInfo {
	var err error
	var storeFileInfos []*model.StoreFileInfo
	if err = dao.DB.Find(&storeFileInfos).Error; err != nil {
		panic(err)
	}
	return storeFileInfos
}

func (s *StoreFileService) ListByStoreDirectoryPath(storeDirectoryPath string) []*model.StoreFileInfo {
	var err error
	var storeFileInfos []*model.StoreFileInfo
	db := dao.DB.Where("`store_directory_path`=?", storeDirectoryPath)
	if err = db.Find(&storeFileInfos).Error; err != nil {
		panic(err)
	}
	return storeFileInfos
}

func (s *StoreFileService) Get(contentHash string) *model.StoreFileInfo {
	var err error
	var storeFileInfo model.StoreFileInfo
	db := dao.DB.Where("`content_hash`=?", contentHash)
	if err = db.Find(&storeFileInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &storeFileInfo
}

// GetAvailableFilename 生成未使用的文件名
func (s *StoreFileService) GetAvailableFilename(userId int, directoryPath, filename string) string {
	fileExtension := filepath.Ext(filename)
	filenameWithoutExtension := filename[0 : len(filename)-len(fileExtension)]
	fileInfos := s.MySpaceService.List(userId, directoryPath, filenameWithoutExtension)
	fileIndex := 0
	var mySpaceFilename string
	for {
		if fileIndex == 0 {
			mySpaceFilename = filenameWithoutExtension + fileExtension
		} else {
			mySpaceFilename = filenameWithoutExtension + fmt.Sprintf("(%d)%s", fileIndex, fileExtension)
		}

		found := false
		for _, fileInfo := range fileInfos {
			if fileInfo.Filename == mySpaceFilename {
				found = true
				break
			}
		}
		if !found {
			break
		}
		fileIndex++
	}
	return mySpaceFilename
}

func (s *StoreFileService) Save(
	userId int,
	contentHash string,
	part *multipart.Part,
	directoryPath string) error {
	// 开事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var err error
	now := time.Now()
	contentType := part.Header.Get("Content-Type") // 内容哈希
	var fileSize int64                             // 文件大小

	storeFileInfo := s.Get(contentHash)
	// 未上传
	if storeFileInfo == nil {
		// 获取文件上传信息，并更新
		ongoingUploadInfo := s.OngoingUploadService.Get(userId, contentHash)
		if ongoingUploadInfo == nil {
			// 不存在则自动创建
			ongoingUploadInfo, err = s.OngoingUploadService.Add(userId, contentHash, contentType)
		} else {
			err = s.OngoingUploadService.Update(userId, contentHash, contentType)
		}
		if err != nil {
			return err
		}

		// 保存文件内容
		tmpFilePath := filepath.Join(ongoingUploadInfo.DirectoryPath, ongoingUploadInfo.Filename)
		tmpFile, err := os.OpenFile(tmpFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = io.Copy(tmpFile, part)
		if err != nil {
			tmpFile.Close()
			return err
		}
		tmpFile.Close()

		// 重新计算内容哈希
		tmpFile, err = os.Open(tmpFilePath)
		if err != nil {
			return err
		}
		md5 := md5.New()
		_, err = io.Copy(md5, tmpFile)
		if err != nil {
			tmpFile.Close()
			return err
		}
		tmpFile.Close()
		contentHash = hex.EncodeToString(md5.Sum(nil))

		// 计算文件大小
		fileSize, err = utils.GetFileSize(filepath.Join(ongoingUploadInfo.DirectoryPath, ongoingUploadInfo.Filename))

		// 保存存储信息
		storeFilename := fmt.Sprintf("%s_%s", now.Format("20060102150405"), part.FileName())
		storeFileInfo = &model.StoreFileInfo{
			ContentHash:        contentHash,
			StoreDirectoryPath: ongoingUploadInfo.DirectoryPath,
			StoreFilename:      storeFilename,
			FileSize:           fileSize,
			MimeType:           contentType,
			CreateTime:         now,
		}
		if err := tx.Create(&storeFileInfo).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 重命名文件
		err = os.Rename(tmpFilePath, filepath.Join(ongoingUploadInfo.DirectoryPath, storeFilename))
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	contentHash = storeFileInfo.ContentHash
	fileSize = storeFileInfo.FileSize
	if err = tx.Where("`user_id`=? AND `content_hash`=?", userId, contentHash).Delete(model.OngoingUploadInfo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加记录到我的空间
	mySpaceFilename := s.GetAvailableFilename(userId, directoryPath, part.FileName())
	fileInfo := model.FileInfo{
		ContentHash:   contentHash,
		UserId:        userId,
		Type:          fileinfotype.Normal,
		DirectoryPath: directoryPath,
		Filename:      mySpaceFilename,
		FileSize:      fileSize,
		MimeType:      contentType,
		UpdateTime:    now,
		CreateTime:    now,
	}
	if err := tx.Create(&fileInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
