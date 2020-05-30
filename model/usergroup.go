package model

import "github.com/farseer810/file-manager/model/constant/memberrole"

type UserGroup struct {
	Group *Group                     `json:"group"`
	User  *User                      `json:"user"`
	Role  memberrole.GroupMemberRole `json:"role"`
}
