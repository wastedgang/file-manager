package main

import (
	"github.com/farseer810/file-manager/cryptoconfig"
	"github.com/farseer810/file-manager/dao"
)

func migrateDatabase() error {
	_, err := cryptoconfig.GetConfiguration()
	if err != nil {
		return nil
	}

	err = dao.InitDatabase()
	if err != nil {
		return err
	}
	return dao.Migrate(dao.DB.DB())
}
