package cryptoconfig

import (
	"fmt"
	"github.com/farseer810/file-manager/model/constant/databasetype"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	AppConfigFilename = "app.config"
	EncryptionKey     = "5pj5d5)sgrh*#$^1lf10d*g*_x6q_*2g9c*=o0#ytv^=b^5gr7"
	SQLiteFilename    = "db.sqlite3"
)

type Config struct {
	Database *DatabaseInfo `json:"database"`
}

type DatabaseInfo struct {
	Type         databasetype.DatabaseType `json:"type"`
	Address      string                    `json:"address"`
	DatabaseName string                    `json:"database_name"`
	Username     string                    `json:"username"`
	Password     string                    `json:"password"`
}

func (d *DatabaseInfo) DBSource() string {
	timeout := 5
	if d.Type == databasetype.MySQL {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&timeout=%ds",
			d.Username, d.Password, d.Address, d.DatabaseName, timeout)
	} else {
		wd, _ := os.Getwd()
		sqlitePath := filepath.Join(wd, SQLiteFilename)
		return fmt.Sprintf("%s?auth&_auth_user=admin&_auth_pass=%s&_auth_crypt=sha1&timeout=%ds",
			sqlitePath, d.Password, timeout)
	}
}

// GetConfiguration 获取配置信息
func GetConfiguration() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(wd, AppConfigFilename)
	encrytedData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config, err := DecryptConfig(encrytedData, []byte(EncryptionKey))
	if err != nil {
		return nil, err
	}
	return config, nil
}

// SaveConfiguration 保存配置信息
func SaveConfiguration(config *Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	configPath := filepath.Join(wd, AppConfigFilename)

	encrypted, err := EncryptConfig(config, []byte(EncryptionKey))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, encrypted, 0755)
	if err != nil {
		return err
	}
	return nil
}
