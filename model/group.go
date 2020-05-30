package model

import (
	"time"
)

type Group struct {
	Id          int       `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	OwnerUserId int       `json:"owner_user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdateTime  time.Time `json:"update_time"`
	CreateTime  time.Time `json:"create_time"`
}

func (Group) TableName() string {
	return "group"
}
