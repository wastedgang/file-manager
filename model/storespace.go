package model

import (
	"time"
)

type StoreSpace struct {
	Id             int       `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	DirectoryPath  string    `json:"directory_path"`
	Remark         string    `json:"remark"`
	AllocateSize   int64     `json:"allocate_size"`
	TotalFileCount int       `json:"total_file_count" gorm:"-"`
	TotalFreeSpace int64     `json:"total_free_space" gorm:"-"`
	CreateTime     time.Time `json:"create_time"`
}

func (StoreSpace) TableName() string {
	return "store_space"
}
