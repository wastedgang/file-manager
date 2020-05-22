package main

import (
	"fmt"
	"github.com/farseer810/file-manager/controller"
	"github.com/farseer810/file-manager/inject"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

type Server interface {
	Run() error
}

type server struct {
	engine *gin.Engine
	Port   int
}

func (s *server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.Port))
}

func GetServer() Server {
	// 解析运行参数
	initFlag()

	// 解析配置文件
	err := ParseConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 初始化日志
	initLog()

	// 迁移数据库
	err = migrateDatabase()
	if err != nil {
		logrus.Fatal(err)
	}

	// 依赖注入
	err = inject.Populate()
	if err != nil {
		panic(err)
	}

	// 初始化web server
	cfg := GetGlobalConfig().ServerConfig
	r := gin.New()
	if cfg.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	err = controller.InitRoutes(r)
	if err != nil {
		logrus.Fatal(err)
	}
	return &server{
		engine: r,
		Port:   cfg.ListenPort,
	}
}
