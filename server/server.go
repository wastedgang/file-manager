package server

import (
	"fmt"
	"github.com/farseer810/file-manager/controller"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// GetServer 初始化server
func GetServer(isDebug bool, listenPort int) Server {
	r := gin.New()
	r.MaxMultipartMemory = 200 << 30 // 200GB上传大小
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	err := controller.InitRoutes(r)
	if err != nil {
		logrus.Fatal(err)
	}
	return &server{
		engine: r,
		Port:   listenPort,
	}
}
