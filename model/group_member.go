package model

import (
	"github.com/farseer810/file-manager/model/constant/memberrole"
	"time"
)

type GroupMember struct {
	Id         int                        `json:"-" gorm:"primary_key,AUTO_INCREMENT"`
	GroupId    int                        `json:"group_id"`
	UserId     int                        `json:"user_id"`
	Role       memberrole.GroupMemberRole `json:"role"`
	UpdateTime time.Time                  `json:"update_time"`
	CreateTime time.Time                  `json:"create_time"`
}

func (GroupMember) TableName() string {
	return "group_member"
}
