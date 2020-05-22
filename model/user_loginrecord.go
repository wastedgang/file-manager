package model

import "time"

type UserLoginRecord struct {
	Id         int       `json:"id" gorm:"primary_key,AUTO_INCREMENT"`
	UserId     int       `json:"user_id"`
	Source     string    `json:"source"`
	CreateTime time.Time `json:"create_time"`
}

func (UserLoginRecord) TableName() string {
	return "user_login_record"
}
