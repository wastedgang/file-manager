package main

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"strings"
)

const (
	DefaultLogFilePrefix = "file-manager"
)

var (
	_globalconfig *GlobalConfig
)

type GlobalConfig struct {
	ServerConfig  *ServerConfig
	LogConfig     *LogConfig
}

type ServerConfig struct {
	ListenPort int
	IsDebug    bool
}

func (s *ServerConfig) ListenAddress() string {
	return fmt.Sprintf(":%d", s.ListenPort)
}

type LogConfig struct {
	Path       string
	Level      string
	FilePrefix string
}

// GetGlobalConfig 获取全局配置信息
func GetGlobalConfig() *GlobalConfig {
	return _globalconfig
}

// ParseServerConfig 解析server配置
func ParseServerConfig(section *ini.Section) (*ServerConfig, error) {
	var err error
	// 解析server配置
	// 处理Port配置
	listenPort := 10001
	if section.HasKey("port") {
		listenPort, err = section.Key("port").Int()
		if err != nil {
			return nil, errors.New(`setting "port" is required in section[server]`)
		}
	}

	serverConfig := &ServerConfig{
		ListenPort: listenPort,
		IsDebug:    isDebug,
	}
	return serverConfig, nil
}

// ParseLogConfig 解析log配置
func ParseLogConfig(section *ini.Section) (*LogConfig, error) {
	// 解析path配置参数
	var path string
	if section.HasKey("path") {
		path = strings.TrimSpace(section.Key("path").String())
	}
	// 计算各种系统下的默认log文件路径
	if path == "" {
		path = "logs"
	}

	// 解析level配置参数
	level := "info"
	if section.HasKey("level") {
		level = strings.TrimSpace(section.Key("level").String())
		switch level {
		case "info":
		case "debug":
		case "warn":
		case "error":
		default:
			return nil, errors.New(`invalid setting "level" in section[log]`)
		}
	}

	// 解析file-prefix配置参数
	filePrefix := DefaultLogFilePrefix
	if section.HasKey("file-prefix") {
		filePrefix = strings.TrimSpace(section.Key("file-prefix").String())
	}
	if filePrefix == "" {
		filePrefix = DefaultLogFilePrefix
	}

	logConfig := &LogConfig{
		Path:       path,
		Level:      level,
		FilePrefix: filePrefix,
	}
	return logConfig, nil
}

// ParseConfig 读取配置文件信息
func ParseConfig() error {
	iniConfig, err := ini.Load(configPath)
	if err != nil {
		return errors.New("invalid config file: failed to read or load config")
	}

	cfg := &GlobalConfig{}
	// 解析server配置
	cfg.ServerConfig, err = ParseServerConfig(iniConfig.Section("server"))
	if err != nil {
		return err
	}

	// 解析日志配置信息
	logSection, err := iniConfig.GetSection("log")
	if err == nil {
		logConfig, err := ParseLogConfig(logSection)
		if err != nil {
			return err
		}
		cfg.LogConfig = logConfig
	}

	_globalconfig = cfg
	return nil
}
