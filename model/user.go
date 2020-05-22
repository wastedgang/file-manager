package model

import (
	"github.com/farseer810/file-manager/model/constant/usertype"
	"time"
)

type User struct {
	Id         int               `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	Type       usertype.UserType `json:"type"`
	Username   string            `json:"username"`
	Password   string            `json:"-"`
	Nickname   string            `json:"nickname"`
	Remark     string            `json:"remark"`
	UpdateTime time.Time         `json:"update_time"`
	CreateTime time.Time         `json:"create_time"`
}

func (User) TableName() string {
	return "user"
}
