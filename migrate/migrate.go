package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	mysqldriver "github.com/golang-migrate/migrate/v4/database/mysql"
	postgresdriver "github.com/golang-migrate/migrate/v4/database/postgres"
	sqlite3driver "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Direction isgo  either up or down.
type DBType string

const (
	MySQL               DBType = "MySQL"
	SQLite3             DBType = "SQLite3"
	Postgres            DBType = "Postgres"
	MigrationsTableName        = "gomigrate_schema_migrations"
)

type MigrationScript struct {
	Version    uint
	Identifier string

	UpScript   string
	DownScript string
}

type Migrate struct {
	migrations       *source.Migrations // 用于检查version和identifier是否重复
	MigrationNames   []string
	MigrationScripts map[string]string
	Versions         []uint
}

func New() *Migrate {
	migrations := &Migrate{
		migrations:       source.NewMigrations(),
		MigrationNames:   make([]string, 0),
		MigrationScripts: make(map[string]string),
		Versions:         make([]uint, 0),
	}
	return migrations
}

// Append 添加迁移版本
func (m *Migrate) Append(version uint, identifier string, upScript string, downScript string) error {
	// 检查version和identifier是否重复
	upMigrationName := fmt.Sprintf("%d_%s.up.sql", version, identifier)
	upMigration := &source.Migration{
		Version:    version,
		Identifier: identifier,
		Direction:  source.Up,
		Raw:        upMigrationName,
	}
	if !m.migrations.Append(upMigration) {
		return errors.New(fmt.Sprintf("duplicate migration(version=%d,identifier:%s)", version, identifier))
	}

	downMigrationName := fmt.Sprintf("%d_%s.down.sql", version, identifier)
	downMigration := &source.Migration{
		Version:    version,
		Identifier: identifier,
		Direction:  source.Down,
		Raw:        downMigrationName,
	}
	if !m.migrations.Append(downMigration) {
		return errors.New(fmt.Sprintf("duplicate migration(version=%d,identifier:%s)", version, identifier))
	}

	// 保存迁移任务信息
	m.MigrationNames = append(m.MigrationNames, upMigrationName)
	m.MigrationScripts[upMigrationName] = upScript
	m.MigrationNames = append(m.MigrationNames, downMigrationName)
	m.MigrationScripts[downMigrationName] = downScript

	m.Versions = append(m.Versions, version)
	return nil
}

// Apply 执行迁移任务
func (m *Migrate) Apply(dbDriverType DBType, db *sql.DB) error {
	if len(m.Versions) == 0 {
		log.Info("migrate succeed!")
		return nil
	}

	// 构造resource和sourceInstance
	resource := bindata.Resource(m.MigrationNames, m.getAssetFunc())
	sourceInstance, err := bindata.WithInstance(resource)
	if err != nil {
		return err
	}

	// 构造databaseDriverName和databaseInstance
	var databaseInstance database.Driver
	var databaseDriverName string
	switch dbDriverType {
	case MySQL:
		databaseDriverName = "mysql"
		databaseInstance, err = mysqldriver.WithInstance(db, &mysqldriver.Config{
			MigrationsTable: MigrationsTableName,
		})
	case SQLite3:
		databaseDriverName = "sqlite3"
		databaseInstance, err = sqlite3driver.WithInstance(db, &sqlite3driver.Config{
			MigrationsTable: MigrationsTableName,
		})
	case Postgres:
		databaseDriverName = "postgres"
		databaseInstance, err = postgresdriver.WithInstance(db, &postgresdriver.Config{
			MigrationsTable: MigrationsTableName,
		})
	default:
		return errors.New(fmt.Sprintf("db %s not supported ", dbDriverType))
	}

	// 构造迁移任务
	migration, err := migrate.NewWithInstance(
		"go-bindata",
		sourceInstance,
		databaseDriverName,
		databaseInstance)
	if err != nil {
		return err
	}

	// 执行迁移

	// hack，将dirty的版本修改成前一个版本号，并将dirty标识去掉
	currentVersion, isDirty, err := migration.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}
	if isDirty {
		if m.Versions[0] == currentVersion {
			// 当前版本是第一个版本时，清空版本迁移表的数据
			_, err = db.Exec(fmt.Sprintf("DELETE FROM %s", MigrationsTableName))
			if err != nil {
				return err
			}
		} else {
			index := 1
			for index < len(m.Versions) {
				if m.Versions[index] == currentVersion {
					break
				}
				index++
			}
			if index >= len(m.Versions) {
				return errors.New(fmt.Sprintf("invalid migration version %d", currentVersion))
			}

			// 当前版本不是第一个版本时，找到上一个版本lastVersion，然后将迁移信息的version、dirty字段设置成lastVersion和0
			lastVersion := m.Versions[index-1]
			err = databaseInstance.SetVersion(int(lastVersion), false)
			if err != nil {
				return err
			}
		}
	}

	err = migration.Up()
	// 没有新迁移不报错
	if err == migrate.ErrNoChange {
		log.Info("migration no change")
		return nil
	}
	if err == nil {
		log.Info("migrate succeed!")
	}
	return err
}

// getAssetFunc 计算所需的AssetFunc
func (m *Migrate) getAssetFunc() bindata.AssetFunc {
	// 闭包
	return func(scripts map[string]string) bindata.AssetFunc {
		return func(name string) ([]byte, error) {
			script, exists := scripts[name]
			if !exists {
				if strings.HasSuffix(name, ".up.sql") {
					name = name[0 : len(name)-len(".up.sql")]
				} else if strings.HasSuffix(name, ".down.sql") {
					name = name[0 : len(name)-len(".down.sql")]
				}
				return nil, errors.New(fmt.Sprintf("migration(%s) not exists", name))
			}
			return []byte(script), nil
		}
	}(m.MigrationScripts)
}