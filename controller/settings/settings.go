package settings

import (
	. "github.com/farseer810/file-manager/controller/vo/statuscode"
	"github.com/farseer810/file-manager/cryptoconfig"
	"github.com/farseer810/file-manager/dao"
	"github.com/farseer810/file-manager/inject"
	"github.com/farseer810/file-manager/model/constant/databasetype"
	"github.com/farseer810/file-manager/service"
	"github.com/gin-gonic/gin"
)

func init() {
	inject.Provide(new(SettingsController))
}

type SettingsController struct {
	SettingsService *service.SettingsService
}

// GetDatabaseSettings 获取数据库配置信息
func (s *SettingsController) GetDatabaseSettings() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		config, err := cryptoconfig.GetConfiguration()
		if err != nil {
			return SystemNotInitialized
		}
		return Success.AddField("settings", config)
	})
	return handler
}

// CheckDatabaseConnection 检查数据库配置信息
func (s *SettingsController) CheckDatabaseSettings() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var form struct {
			Type     string `form:"type" binding:"required,oneof=MySQL SQLite"`
			Address  string `form:"address"`
			Database string `form:"database"`
			Username string `form:"username"`
			Password string `form:"password" binding:"required"`
		}
		if err := ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}

		ok := s.SettingsService.CheckDatabaseConnection(
			databasetype.DatabaseType(form.Type),
			form.Address,
			form.Username,
			form.Password)
		return Success.AddField("ok", ok)
	})
	return handler
}

// SaveDatabaseConfig 保存数据库配置信息
func (s *SettingsController) SetDatabaseSettings() gin.HandlerFunc {
	handler := ConvertGinHandlerFunc(func(ctx *gin.Context) *Response {
		var err error
		var form struct {
			Type     string `form:"type" binding:"required,oneof=MySQL SQLite"`
			Address  string `form:"address"`
			Database string `form:"database"`
			Username string `form:"username"`
			Password string `form:"password" binding:"required"`
		}
		if err = ctx.ShouldBind(&form); err != nil {
			return BadRequest
		}
		if form.Type == "SQLite" {
			form.Address = ""
			form.Database = ""
			form.Username = ""
		}
		// 检查配置信息是否已存在，已存在说明系统已初始化，不能再初始化了
		_, err = cryptoconfig.GetConfiguration()
		if err == nil {
			return SystemInitialized
		}

		databaseInfo := &cryptoconfig.DatabaseInfo{
			Type:         databasetype.DatabaseType(form.Type),
			Address:      form.Address,
			DatabaseName: form.Database,
			Username:     form.Username,
			Password:     form.Password,
		}

		// 测试连接数据库
		ok := s.SettingsService.CheckDatabaseConnection(
			databaseInfo.Type,
			databaseInfo.Address,
			databaseInfo.Username,
			databaseInfo.Password)
		if !ok {
			return DatabaseConnectFail
		}

		// 创建数据库
		err = s.SettingsService.CreateDatabase(databaseInfo)
		if err != nil {
			return DatabaseCreateFail
		}

		// 保存配置信息
		config := &cryptoconfig.Config{Database: databaseInfo}
		err = cryptoconfig.SaveConfiguration(config)
		if err != nil {
			return InternalServerError
		}

		// 迁移sql
		err = dao.InitDatabase(false)
		if err != nil {
			return DatabaseConnectFail
		}
		err = dao.Migrate(dao.DB.DB())
		if err != nil {
			return DatabaseMigrateFail
		}
		return Success
	})
	return handler
}
