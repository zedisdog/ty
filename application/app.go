package application

import (
	"embed"
	"fmt"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/database"
	"github.com/zedisdog/ty/database/migrate"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/log/zap"
	"github.com/zedisdog/ty/sdk/net/http/server"
	"github.com/zedisdog/ty/sdk/net/http/server/gin"
	"github.com/zedisdog/ty/strings"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var instance *App
var once sync.Once

// GetInstance gets the application singleton instance.
func GetInstance() *App {
	once.Do(func() {
		instance = &App{
			httpServers: new(sync.Map),
			databases:   new(sync.Map),
			migrates:    migrate.NewFsDriver(),
		}
	})

	return instance
}

type IApplication interface {
	Init(config config.IConfig)
	Boot()
	RegisterModule(module IModule)
	RegisterHttpServerRoute(f func(serverEngine interface{}) error)
	RegisterMigrate(fs *embed.FS)
	RegisterDatabase(name string, db interface{})
	Database(name string) interface{}
	Run()
	Stop()
	Wait(closeFunc ...func())
}

type App struct {
	Config      config.IConfig
	httpServers *sync.Map
	Logger      log.ILog
	modules     []IModule
	databases   *sync.Map
	migrates    *migrate.EmbedDriver
}

// Init set config to application.
func Init(config config.IConfig) {
	GetInstance().Init(config)
}
func (app *App) Init(config config.IConfig) {
	app.Config = config

	if logConfig := app.Config.Sub("log"); logConfig != nil {
		app.initLog(app.Config.Sub("log"))
	} else {
		panic(errx.New("[application] there is no log config"))
	}
}

func (app *App) initLog(config config.IConfig) {
	driver := config.GetString("driver", "zap")
	fmt.Printf("[application] init log using %s...\n", driver)
	switch driver {
	case "zap":
		app.Logger = zap.NewZapLog()
	}
	return
}

func (app *App) RegisterDatabase(name string, db interface{}) {
	_, exists := app.databases.Load(name)
	if exists {
		panic(errx.New(fmt.Sprintf("database <%s> is already exists", name)))
	}
	app.databases.Store(name, db)
}

func (app *App) migrate() {
	var (
		migrator migrate.IMigrator = &migrate.DefaultMigrator{}
		err      error
	)

	if err != nil {
		panic(err)
	}

	err = migrator.Migrate(strings.EncodeQuery(app.Config.GetString("default.database.dsn")), app.migrates)
	if err != nil {
		panic(err)
	}
}

// Boot boots the application.
func Boot() {
	GetInstance().Boot()
}
func (app *App) Boot() {
	app.migrate()
}

func (app *App) newDB(config map[string]interface{}) (db *gorm.DB, err error) {
	return database.NewDatabase(strings.EncodeQuery(config["dsn"].(string)))
}

// RegisterModule register module to application.
func RegisterModule(module IModule) {
	GetInstance().RegisterModule(module)
}
func (app *App) RegisterModule(module IModule) {
	app.Logger.Info("[application] register module", &log.Field{Name: "name", Value: module.Name()})
	err := module.Register(app)
	if err != nil {
		panic(errx.Wrap(err, "[application] register module failed"))
	}
	app.modules = append(app.modules, module)
}

// RegisterHttpServerRoute register module routes to application.
func RegisterHttpServerRoute(f func(serverEngine interface{}) error) {
	GetInstance().RegisterHttpServerRoute(f)
}
func (app *App) RegisterHttpServerRoute(f func(serverEngine interface{}) error) {
	svr, ok := app.httpServers.Load("default")
	if !ok {
		def := app.Config.Sub("default.httpServer")
		if def != nil && def.GetBool("enable") {
			app.Logger.Info("[application] create default http server...")
			svr = gin.NewGinServer(fmt.Sprintf(
				"%s:%d",
				def.GetString("host"),
				def.GetInt("port"),
			))
			app.httpServers.Store("default", svr)
		} else {
			app.Logger.Info("[application] no default http server specifies")
			return
		}
	}

	err := errx.Wrap(svr.(server.IHTTPServer).RegisterRoutes(f), "register route failed")
	if err != nil {
		panic(err)
	}
}

func Run() {
	GetInstance().Run()
}
func (app *App) Run() {
	app.httpServers.Range(func(key, value any) bool {
		app.Logger.Info(
			"[application] run http server...",
			&log.Field{Name: "name", Value: key},
		)
		value.(server.IHTTPServer).Run()
		return true
	})
}

func Stop() {
	GetInstance().Stop()
}
func (app *App) Stop() {
	app.httpServers.Range(func(key, value any) bool {
		app.Logger.Info(
			"[application] shutdown http server...",
			&log.Field{Name: "name", Value: key},
		)
		err := value.(server.IHTTPServer).Shutdown()
		if err != nil {
			app.Logger.Error(
				"[application]shutdown http server error",
				&log.Field{Name: "name", Value: key},
				&log.Field{Name: "error", Value: err},
			)
		}
		return true
	})
}

func Wait(closeFunc ...func()) {
	GetInstance().Wait(closeFunc...)
}
func (app *App) Wait(closeFunc ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	for _, cls := range closeFunc {
		cls()
	}
}

func RegisterMigrate(fs *embed.FS) {
	GetInstance().migrates.Add(fs)
}
func (app *App) RegisterMigrate(fs *embed.FS) {
	app.migrates.Add(fs)
}

func Database(name string) interface{} {
	return GetInstance().Database(name)
}
func (app *App) Database(name string) (db interface{}) {
	db, _ = app.databases.Load(name)
	return
}

//func (app *App) bootModules(config config.IConfig) {
//	for _, module := range app.modules {
//		if config.IsSet(module.Name()) {
//			err := module.Boot(config.Sub(module.Name()))
//			if err != nil {
//				panic(errx.Wrap(err, "[app]boot module failed"))
//			}
//		}
//	}
//}
//
//func RegisterHttpServer(name string, server server.IHTTPServer) {
//	instance.RegisterHttpServer(name, server)
//}
//func (app *App) RegisterHttpServer(name string, server server.IHTTPServer) {
//	app.httpServers.Store(name, server)
//}
//
//func GetHttpServer(name string) server.IHTTPServer {
//	return instance.GetHttpServer(name)
//}
//func (app *App) GetHttpServer(name string) server.IHTTPServer {
//	v, exists := app.httpServers.Load(name)
//	if !exists {
//		return nil
//	}
//
//	return v.(server.IHTTPServer)
//}
