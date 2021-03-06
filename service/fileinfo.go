package service

import (
	"fmt"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/farseer810/file-manager/model/constant/fileactiontype"
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/jinzhu/gorm"
	"path/filepath"
	"strings"
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

// GetAvailableFilename 生成未使用的文件名
func (s *FileInfoService) GetAvailableFilename(userId int, directoryPath, filename string) string {
	fileExtension := filepath.Ext(filename)
	filenameWithoutExtension := filename[0 : len(filename)-len(fileExtension)]
	fileInfos := s.List(userId, directoryPath, filenameWithoutExtension)
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
	updates := map[string]interface{}{"filename": newFilename, "update_time": now}
	if err = tx.Model(&model.FileInfo{Id: oldFileInfo.Id}).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}
	if oldFileInfo.Type == fileinfotype.Directory {
		// 处理子文件夹
		oldDirectoryPath := filepath.Join(oldFileInfo.DirectoryPath, oldFileInfo.Filename)
		newDirectoryPath := filepath.Join(oldFileInfo.DirectoryPath, newFilename)
		var subFiles []*model.FileInfo
		directoryPathParam := fmt.Sprintf("%s/%%", oldDirectoryPath)
		db := tx.Where("`user_id`=? AND (`directory_path`=? OR `directory_path` LIKE ?)", oldFileInfo.UserId, oldDirectoryPath, directoryPathParam)
		if err = db.Find(&subFiles).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, subFile := range subFiles {
			newSubDirectoryPath := strings.Replace(subFile.DirectoryPath, oldDirectoryPath, newDirectoryPath, 1)
			if err = tx.Model(&subFile).Update("directory_path", newSubDirectoryPath).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *FileInfoService) ListByFilenames(userId int, directoryPath string, filenames []string) []*model.FileInfo {
	var fileInfos []*model.FileInfo
	db := dao.DB.Where("`user_id`=? AND `directory_path`=? AND filename IN (?)", userId, directoryPath, filenames)
	if err := db.Find(&fileInfos).Error; err != nil {
		panic(err)
	}
	return fileInfos
}

func (s *FileInfoService) ListSubFiles(userId int, directoryPath string) []*model.FileInfo {
	var fileInfos []*model.FileInfo
	directoryPathParam := fmt.Sprintf("%s/%%", directoryPath)
	db := dao.DB.Where("`user_id`=? AND (`directory_path`=? OR `directory_path` LIKE ?)", userId, directoryPath, directoryPathParam)
	if err := db.Find(&fileInfos).Error; err != nil {
		panic(err)
	}
	return fileInfos
}

func (s *FileInfoService) Delete(userId int, directoryPath string, filenames []string) error {
	var err error
	// 开事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	fileInfos := s.ListByFilenames(userId, directoryPath, filenames)
	// 删除指定文件
	db := tx.Where("`user_id`=? AND `directory_path`=? AND filename IN (?)", userId, directoryPath, filenames)
	if err = db.Delete(model.FileInfo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除子文件夹
	for _, fileInfo := range fileInfos {
		if fileInfo.Type != fileinfotype.Directory {
			continue
		}
		subDirectoryPath := filepath.Join(fileInfo.DirectoryPath, fileInfo.Filename)
		directoryPathParam := fmt.Sprintf("%s/%%", subDirectoryPath)
		db := tx.Where("`user_id`=? AND (`directory_path`=? OR `directory_path` LIKE ?)", userId, subDirectoryPath, directoryPathParam)
		if err = db.Delete(model.FileInfo{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *FileInfoService) Copy(sourceFileInfo *model.FileInfo, targetDirectoryPath string, actionType fileactiontype.FileActionType) error {
	if sourceFileInfo.DirectoryPath == targetDirectoryPath && actionType == fileactiontype.Override {
		return nil
	}

	userId := sourceFileInfo.UserId
	var err error
	// 开事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	targetFilename := sourceFileInfo.Filename
	// 删除需要覆盖的文件以及其子文件
	if actionType == fileactiontype.Override {
		// 删除需要覆盖的文件
		db := tx.Where("`user_id`=? AND `directory_path`=? AND `filename`=?", userId, targetDirectoryPath, sourceFileInfo.Filename)
		if err = db.Delete(model.FileInfo{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 如果有子文件，也一并删除
		directoryPathParam := fmt.Sprintf("%s/%s/%%", targetDirectoryPath, sourceFileInfo.Filename)
		db = tx.Where("`user_id`=? AND `directory_path` LIKE ?", userId, directoryPathParam)
		if err = db.Delete(model.FileInfo{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 计算新文件名
		targetFilename = s.GetAvailableFilename(userId, targetDirectoryPath, sourceFileInfo.Filename)
	}

	now := time.Now()
	// 先复制指定文件
	targetFileInfo := &model.FileInfo{
		ContentHash:   sourceFileInfo.ContentHash,
		UserId:        sourceFileInfo.UserId,
		Type:          sourceFileInfo.Type,
		DirectoryPath: targetDirectoryPath,
		Filename:      targetFilename,
		FileSize:      sourceFileInfo.FileSize,
		MimeType:      sourceFileInfo.MimeType,
		UpdateTime:    now,
		CreateTime:    now,
	}
	if err := dao.DB.Create(&targetFileInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if sourceFileInfo.Type != fileinfotype.Directory {
		return tx.Commit().Error
	}

	// 若为文件夹，则该文件夹下的所有子文件夹与子文件均需要复制
	sourceFilePath := filepath.Join(sourceFileInfo.DirectoryPath, sourceFileInfo.Filename)
	subFileInfos := s.ListSubFiles(userId, sourceFilePath)
	if len(subFileInfos) == 0 {
		return tx.Commit().Error
	}

	for _, fileInfo := range subFileInfos {
		subFileDirectoryPath := filepath.Join(targetDirectoryPath, targetFilename)
		newFileInfo := &model.FileInfo{
			ContentHash:   fileInfo.ContentHash,
			UserId:        fileInfo.UserId,
			Type:          fileInfo.Type,
			DirectoryPath: strings.Replace(fileInfo.DirectoryPath, sourceFilePath, subFileDirectoryPath, 1),
			Filename:      fileInfo.Filename,
			FileSize:      fileInfo.FileSize,
			MimeType:      fileInfo.MimeType,
			UpdateTime:    now,
			CreateTime:    now,
		}
		if err := dao.DB.Create(&newFileInfo).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (s *FileInfoService) Move(userId int, oldDirectoryPath string, filenames []string, newDirectoryPath string) error {
	// 原地移动即不需要移动
	if oldDirectoryPath == newDirectoryPath || len(filenames) == 0 {
		return nil
	}

	// 开事务
	tx := dao.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 查询欲移动文件的信息列表
	fileInfos := s.ListByFilenames(userId, oldDirectoryPath, filenames)
	if len(fileInfos) == 0 {
		return nil
	}
	// 若为文件夹，则该文件夹下的所有子文件夹与子文件均需要移动
	for _, fileInfo := range fileInfos {
		if fileInfo.Type != fileinfotype.Directory {
			continue
		}
		subFileInfos := s.ListSubFiles(userId, filepath.Join(fileInfo.DirectoryPath, fileInfo.Filename))
		if len(subFileInfos) == 0 {
			continue
		}
		fileInfos = append(fileInfos, subFileInfos...)
	}

	now := time.Now()
	for _, fileInfo := range fileInfos {
		fileInfo.DirectoryPath = strings.Replace(fileInfo.DirectoryPath, oldDirectoryPath, newDirectoryPath, 1)
		fileInfo.UpdateTime = now
		tx.Save(&fileInfo)
	}
	return tx.Commit().Error
}
