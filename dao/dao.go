package dao

import (
	"database/sql"
	"errors"
	"github.com/farseer810/file-manager/cryptoconfig"
	"github.com/farseer810/file-manager/inject"
	mysqlmigrate "github.com/farseer810/file-manager/migrate/mysql"
	sqlite3migrate "github.com/farseer810/file-manager/migrate/sqlite3"
	"github.com/farseer810/file-manager/model/constant/databasetype"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrInvalidDatabaseType = errors.New("invalid database type")
	DB                     *gorm.DB
)

func InitDatabase() error {
	if DB != nil {
		return nil
	}
	config, err := cryptoconfig.GetConfiguration()
	if err != nil {
		return err
	}
	
	switch config.Database.Type {
	case databasetype.MySQL:
		DB, err = gorm.Open("mysql", config.Database.DBSource())
	case databasetype.SQLite:
		DB, err = gorm.Open("sqlite3", config.Database.DBSource())
	default:
		return ErrInvalidDatabaseType
	}
	if err != nil {
		return err
	}

	inject.Provide(DB)
	return nil
}

func Migrate(db *sql.DB) error {
	config, err := cryptoconfig.GetConfiguration()
	if err != nil {
		return err
	}

	if config.Database.Type == databasetype.MySQL {
		err = mysqlmigrate.Migrate(db)
	} else {
		err = sqlite3migrate.Migrate(db)
	}
	return err
}
