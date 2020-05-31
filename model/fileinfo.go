package model

import (
	"github.com/farseer810/file-manager/model/constant/fileinfotype"
	"github.com/farseer810/file-manager/utils"
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

func (f *FileInfo) Less(f2 *FileInfo) bool {
	if f.Type == fileinfotype.Directory && f2.Type != fileinfotype.Directory {
		return true
	}
	if f.Type != fileinfotype.Directory && f2.Type == fileinfotype.Directory {
		return false
	}
	filename1 := []rune(f.Filename)
	filename2 := []rune(f2.Filename)
	length1 := len(filename1)
	length2 := len(filename2)
	minLength := length1
	if length2 < length1 {
		minLength = length2
	}

	// 顺序依次为：其他字符、数字、字母、汉字
	var r1, r2 rune
	var pinyin1, pinyin2 string
	for i := 0; i < minLength; i++ {
		r1, r2 = filename1[i], filename2[i]
		// 字符一样，跳过，比较下一个
		if r1 == r2 {
			continue
		}
		pinyin1 = utils.GetPinYin(r1)
		pinyin2 = utils.GetPinYin(r2)
		if pinyin1 == "" && pinyin2 == "" {
			// 两个字符都不是汉字
			isLetter1 := utils.IsLetter(r1)
			isLetter2 := utils.IsLetter(r2)
			if isLetter1 && !isLetter2 {
				return false
			} else if !isLetter1 && isLetter2 {
				return true
			} else if !isLetter1 && !isLetter2 {
				// 两个字符虽不是汉字，又不是字母
				isDigit1 := utils.IsDigit(r1)
				isDigit2 := utils.IsDigit(r2)
				if isDigit1 && !isDigit2 {
					return false
				} else if !isDigit1 && isDigit2 {
					return true
				}
			}
			return r1 < r2
		} else if pinyin1 != "" && pinyin2 == "" {
			return false
		} else if pinyin1 == "" && pinyin2 != "" {
			return true
		}
		return pinyin1 < pinyin2
	}
	return length1 < length2
}
