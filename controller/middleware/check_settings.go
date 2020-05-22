package middleware

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/cryptoconfig"
	"github.com/farseer810/file-manager/dao"
	"github.com/gin-gonic/gin"
)

func CheckConfigurationHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据库已初始化
		if dao.DB != nil {
			ctx.Next()
			return
		}

		// 以下两个接口跳过检验
		if ctx.FullPath() == "/api/v1/settings/database" || ctx.FullPath() == "/api/v1/settings/database/check" {
			ctx.Next()
			return
		}

		_, err := cryptoconfig.GetConfiguration()
		if err != nil {
			ctx.JSON(200, SystemNotInitialized)
			return
		}
		err = dao.InitDatabase()
		if err != nil || dao.DB == nil {
			ctx.JSON(200, DatabaseConnectFail)
			return
		}
		ctx.Next()
	}
}
