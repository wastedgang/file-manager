package model

import (
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"time"
)

type FileInfo struct {
	Id            int                       `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	ContentHash   string                    `json:"content_hash"`
	UserId        int                       `json:"user_id"`
	Type          fileinfotype.FileInfoType `json:"type"`
	DirectoryPath string                    `json:"directory_path"`
	Filename      string                    `json:"filename"`
	FileSize      int64                     `json:"file_size"`
	MimeType      string                    `json:"mime_type"`
	UpdateTime    time.Time                 `json:"update_time"`
	CreateTime    time.Time                 `json:"create_time"`
}

func (FileInfo) TableName() string {
	return "file_info"
}
