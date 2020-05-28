package service

import (
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	inject.Provide(new(StoreSpaceService))
}

type StoreSpaceService struct {
	StoreFileService *StoreFileService
}

func (s *StoreSpaceService) Add(directoryPath string, allocateSize int64, remark string) error {
	now := time.Now()
	storeSpace := model.StoreSpace{
		DirectoryPath: directoryPath,
		AllocateSize:  allocateSize,
		Remark:        remark,
		CreateTime:    now,
	}
	if err := dao.DB.Create(&storeSpace).Error; err != nil {
		return err
	}
	return nil
}

func (s *StoreSpaceService) Delete(directoryPath string) error {
	// TODO: 迁移数据

	if err := dao.DB.Where("`directory_path`=?", directoryPath).Delete(model.StoreSpace{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *StoreSpaceService) Update(directoryPath string, allocateSize int64, remark string) error {
	updates := map[string]interface{}{
		"allocate_size": allocateSize,
		"remark":        remark,
	}
	if err := dao.DB.Model(model.User{}).Where("`directory_path`=?", directoryPath).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (s *StoreSpaceService) GetByDirectoryPath(directoryPath string) *model.StoreSpace {
	var storeSpace model.StoreSpace
	if err := dao.DB.First(&storeSpace, "`directory_path`=?", directoryPath).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		panic(err)
	}
	return &storeSpace
}

func (s *StoreSpaceService) List() []*model.StoreSpace {
	var err error
	var storeSpaces []*model.StoreSpace
	if err = dao.DB.Find(&storeSpaces).Error; err != nil {
		panic(err)
	}

	// 计算空间文件数量，以及空间文件总大小
	storeFileCountMap := make(map[string]int)
	storeFileSizeMap := make(map[string]int64)
	for _, storeFileInfo := range s.StoreFileService.List() {
		storeDirectoryPath := storeFileInfo.StoreDirectoryPath
		storeFileCountMap[storeDirectoryPath] = storeFileCountMap[storeDirectoryPath] + 1
		storeFileSizeMap[storeDirectoryPath] = storeFileSizeMap[storeDirectoryPath] + storeFileInfo.FileSize
	}

	for _, storeSpace := range storeSpaces {
		storeSpace.TotalFileCount = storeFileCountMap[storeSpace.DirectoryPath]
		// 计算空间文件剩余空间
		storeSpace.TotalFreeSpace = storeSpace.AllocateSize - storeFileSizeMap[storeSpace.DirectoryPath]
	}
	return storeSpaces
}

// GetBestStoreSpace 计算剩余空间最大的存储空间
func (s *StoreSpaceService) GetBestStoreSpace() *model.StoreSpace {
	storeSpaces := s.List()
	if len(storeSpaces) == 0 {
		return nil
	}
	bestStoreSpace := storeSpaces[0]
	for i := 1; i < len(storeSpaces); i++ {
		if bestStoreSpace.TotalFreeSpace < storeSpaces[i].TotalFreeSpace {
			bestStoreSpace = storeSpaces[i]
		}
	}
	return bestStoreSpace
}
