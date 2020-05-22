package main

import (
	"fmt"
	"github.com/farseer810/file-manager/dao"
)

func migrateDatabase() error {
	err := dao.InitDatabase()
	if err != nil || dao.DB == nil {
		fmt.Println(err)
		return nil
	}

	return dao.Migrate(dao.DB.DB())
}