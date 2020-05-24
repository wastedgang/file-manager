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

type StoreSpaceService struct{}

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

	// TODO: 计算空间文件数量
	// TODO: 计算空间文件剩余空间
	return storeSpaces
}
