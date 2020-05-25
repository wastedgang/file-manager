package service

import (
	"fmt"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
	"time"
)

func init() {
	inject.Provide(new(OngoingUploadService))
}

type OngoingUploadService struct {
	StoreSpaceService *StoreSpaceService
}

// GetUploadStartPoint 获取正在上传的文件大小
func (o *OngoingUploadService) GetUploadStartPoint(userId int, contentHash string) int64 {
	ongoingUploadInfo := o.Get(userId, contentHash)
	if ongoingUploadInfo == nil {
		return 0
	}

	fileInfo, err := os.Stat(filepath.Join(ongoingUploadInfo.DirectoryPath, ongoingUploadInfo.Filename))
	if err != nil || fileInfo == nil {
		return 0
	}
	return fileInfo.Size()
}

// Get 获取正在上传文件的信息
func (o *OngoingUploadService) Get(userId int, contentHash string) *model.OngoingUploadInfo {
	var err error
	var ongoingUploadInfo model.OngoingUploadInfo
	db := dao.DB.Where("`user_id`=? AND `content_hash`=?", userId, contentHash)
	if err = db.Find(&ongoingUploadInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &ongoingUploadInfo
}

// Add 创建当前上传记录
func (o *OngoingUploadService) Add(userId int, contentHash, contentType string) (*model.OngoingUploadInfo, error) {
	// 计算存储位置
	bestStoreSpace := o.StoreSpaceService.GetBestStoreSpace()
	filename := fmt.Sprintf("%d_%s.upload.tmp", userId, contentHash)

	now := time.Now()
	ongoingUploadInfo := model.OngoingUploadInfo{
		ContentHash:   contentHash,
		UserId:        userId,
		DirectoryPath: bestStoreSpace.DirectoryPath,
		Filename:      filename,
		MimeType:      contentType,
		CreateTime:    now,
	}
	if err := dao.DB.Create(&ongoingUploadInfo).Error; err != nil {
		panic(err)
	}
	return &ongoingUploadInfo, nil
}

// Update 更新当前上传记录
func (o *OngoingUploadService) Update(userId int, contentHash, contentType string) error {
	db := dao.DB.Where("`user_id`=? AND `content_hash`=?", userId, contentHash)

	// 更新mime_type
	updates := map[string]interface{}{"mime_type": contentType}
	if err := db.Model(model.OngoingUploadInfo{}).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
