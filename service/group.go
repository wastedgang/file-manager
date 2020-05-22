package service

import "github.com/farseer810/file-manager/inject"

func init() {
	inject.Provide(new(GroupService))
}

type GroupService struct {}
