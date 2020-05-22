package service

import (
	"database/sql"
	"fmt"
	"github.com/farseer810/file-manager/cryptoconfig"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/databasetype"
)

func init() {
	inject.Provide(new(SettingsService))
}

type SettingsService struct{}

// getDB 计算相应类型的*sql.DB
func (s *SettingsService) getDB(databaseInfo *cryptoconfig.DatabaseInfo) (*sql.DB, error) {
	var driverName string
	if databasetype.MySQL == databaseInfo.Type {
		driverName = "mysql"
	} else {
		driverName = "sqlite3"
	}
	dataSourceName := databaseInfo.DBSource()

	return sql.Open(driverName, dataSourceName)
}

// CheckDatabaseConnection 测试连接数据库
func (s *SettingsService) CheckDatabaseConnection(databaseType databasetype.DatabaseType, address, username, password string) bool {
	databaseInfo := &cryptoconfig.DatabaseInfo{
		Type:         databaseType,
		Address:      address,
		DatabaseName: "",
		Username:     username,
		Password:     password,
	}
	db, err := s.getDB(databaseInfo)
	if err != nil {
		return false
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return false
	}
	return true
}

func (s *SettingsService) CreateDatabase(databaseInfo *cryptoconfig.DatabaseInfo) error {
	db, err := s.getDB(&cryptoconfig.DatabaseInfo{
		Type:         databaseInfo.Type,
		Address:      databaseInfo.Address,
		DatabaseName: "",
		Username:     databaseInfo.Username,
		Password:     databaseInfo.Password,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	if databaseInfo.Type == databasetype.MySQL {
		// 检查数据库是否存在
		rows, err := db.Query("show databases")
		if err != nil {
			return err
		}
		defer rows.Close()
		var databaseName string
		databaseNames := make([]string, 0)
		for rows.Next() {
			err := rows.Scan(&databaseName)
			if err != nil {
				return err
			}
			databaseNames = append(databaseNames, databaseName)
		}

		// 数据库是否存在
		for _, name := range databaseNames {
			if name == databaseInfo.DatabaseName {
				return nil
			}
		}

		// 不存在数据库，创建新数据库
		sqlScript := fmt.Sprintf("create database `%s` default character set utf8 collate utf8_general_ci", databaseInfo.DatabaseName)
		_, err = db.Exec(sqlScript)
		if err != nil {
			return err
		}
	}
	return nil
}
