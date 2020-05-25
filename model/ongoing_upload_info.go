package model

import (
	"time"
)

type OngoingUploadInfo struct {
	Id            int                       `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	ContentHash   string                    `json:"content_hash"`
	UserId        int                       `json:"user_id"`
	DirectoryPath string                    `json:"directory_path"`
	Filename      string                    `json:"filename"`
	MimeType      string                    `json:"mime_type"`
	CreateTime    time.Time                 `json:"create_time"`
}

func (OngoingUploadInfo) TableName() string {
	return "ongoing_upload_info"
}
