package main

import (
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// 解析运行参数
	initFlag()

	// 解析配置文件
	err := ParseConfig()
	if err != nil {
		os.Exit(1)
	}

	// 初始化日志
	initLog()

	// 依赖注入
	err = inject.Populate()
	if err != nil {
		panic(err)
	}

	// 迁移数据库
	err = migrateDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	// 运行服务器
	serverConfig := GetGlobalConfig().ServerConfig
	server := server.GetServer(serverConfig.IsDebug, serverConfig.ListenPort)
	server.Run()
}
