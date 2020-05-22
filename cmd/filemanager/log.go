package main

import (
	mylog "github.com/farseer810/file-manager/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	formatter := &mylog.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	}
	log.SetFormatter(formatter)
}

// initLog 初始化日志
func initLog() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	formatter := &mylog.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	}
	log.SetFormatter(formatter)

	// 调试模式
	if GetGlobalConfig().ServerConfig.IsDebug {
		log.SetLevel(log.TraceLevel)
	} else {
		cfg := GetGlobalConfig().LogConfig
		if cfg == nil {
			return
		}

		// 设置log级别
		level, err := log.ParseLevel(cfg.Level)
		if err != nil {
			log.Error("invalid log level, fallback to use info")
		}
		log.SetLevel(level)

		// 滚动日志
		rotatePattern := strings.TrimRight(cfg.Path, "/") + "/"
		if cfg.FilePrefix != "" {
			rotatePattern += cfg.FilePrefix + "-"
		}
		rotatePattern += "%Y-%m-%d.log"
		rl, err := rotatelogs.New(rotatePattern)
		if err != nil {
			log.Fatal("failed to create logs:", err)
			return
		}
		log.SetOutput(io.MultiWriter(os.Stdout, rl))
	}
}
