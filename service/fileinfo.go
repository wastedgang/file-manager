package service

import (
	"fmt"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/jinzhu/gorm"
	"path/filepath"
	"time"
)

func init() {
	inject.Provide(new(FileInfoService))
}

type FileInfoService struct {
}

func (s *FileInfoService) ListDirectories(userId int) []*model.FileInfo {
	var err error
	var fileInfos []*model.FileInfo
	if err = dao.DB.Where("`user_id`=? AND `type`=?", userId, fileinfotype.Directory).Find(&fileInfos).Error; err != nil {
		panic(err)
	}
	return fileInfos
}

func (s *FileInfoService) List(userId int, directoryPath string, searchWord string) []*model.FileInfo {
	var err error
	var fileInfos []*model.FileInfo

	var db *gorm.DB
	if searchWord == "" {
		db = dao.DB.Where("`user_id`=? AND `directory_path`=?", userId, directoryPath)
	} else {
		searchWordParam := fmt.Sprintf("%%%s%%", searchWord)
		db = dao.DB.Where("`user_id`=? AND `directory_path`=? AND `filename` LIKE ?", userId, directoryPath, searchWordParam)
	}
	if err = db.Find(&fileInfos).Error; err != nil {
		panic(err)
	}
	return fileInfos
}

func (s *FileInfoService) IsDirectoryExists(userId int, path string) bool {
	// 忽略根目录
	if path == "/" {
		return true
	}

	fileInfo := s.Get(userId, path)
	return fileInfo != nil && fileInfo.Type == fileinfotype.Directory
}

func (s *FileInfoService) Get(userId int, path string) *model.FileInfo {
	// 计算所在目录路径和文件名
	var err error
	path, err = filepath.Abs(filepath.Clean(path))
	if err != nil {
		panic(err)
	}
	basename := filepath.Base(path)
	direname := filepath.Dir(path)

	var fileInfo model.FileInfo
	db := dao.DB.Where("`user_id`=? AND `directory_path`=? AND `filename`=?", userId, direname, basename)
	if err = db.Find(&fileInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &fileInfo
}

func (s *FileInfoService) CreateDirectory(userId int, path string) error {
	// 忽略根目录
	var err error
	path, err = filepath.Abs(filepath.Clean(path))
	if err != nil {
		panic(err)
	}
	if path == "/" {
		return nil
	}

	// 计算所在目录路径和文件名
	direname := filepath.Dir(path)
	basename := filepath.Base(path)

	// 创建文件夹记录
	now := time.Now()
	fileInfo := model.FileInfo{
		UserId:        userId,
		Type:          fileinfotype.Directory,
		DirectoryPath: direname,
		Filename:      basename,
		FileSize:      0,
		UpdateTime:    now,
		CreateTime:    now,
	}
	if err := dao.DB.Create(&fileInfo).Error; err != nil {
		return err
	}
	return nil
}

func (s *FileInfoService) Rename(oldFileInfo *model.FileInfo, newFilename string) error {
	var err error
	// 开事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	now := time.Now()
	if err = tx.Model(oldFileInfo).Updates(map[string]interface{}{"filename": newFilename, "update_time": now}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if oldFileInfo.Type == fileinfotype.Directory {
		// TODO: 处理子文件夹
	}

	return tx.Commit().Error
}
