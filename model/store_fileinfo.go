package model

import (
	"time"
)

type StoreFileInfo struct {
	Id                 int       `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	ContentHash        string    `json:"content_hash"`
	StoreDirectoryPath string    `json:"store_directory_path"`
	StoreFilename      string    `json:"store_filename"`
	FileSize           int64     `json:"file_size"`
	MimeType           string    `json:"mime_type"`
	UpdateTime         time.Time `json:"update_time"`
	CreateTime         time.Time `json:"create_time"`
}

func (StoreFileInfo) TableName() string {
	return "store_file_info"
}
