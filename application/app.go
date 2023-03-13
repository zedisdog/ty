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
	"github.com/zedisdog/ty/storage"
	"github.com/zedisdog/ty/storage/drivers"
	"github.com/zedisdog/ty/strings"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"reflect"
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
			modules:     new(sync.Map),
			components:  new(sync.Map),
			migrates:    migrate.NewFsDriver(),
		}
	})

	return instance
}

type IApplication interface {
	Init(config config.IConfig)
	Boot()
	Run()
	Stop()
	Wait(closeFunc ...func())

	RegisterModule(module interface{})
	RegisterHttpServerRoute(f func(serverEngine interface{}) error)
	RegisterMigrate(fs *embed.FS)
	RegisterDatabase(name string, db interface{})
	Database(name string) interface{}
	Storage() storage.IStorage
	Logger() log.ILog
	Module(nameOrType interface{}) (module interface{})
	Config() config.IConfig

	SetComponent(key any, value any)
	GetComponent(key any) any
}

type App struct {
	config      config.IConfig
	httpServers *sync.Map
	logger      log.ILog
	modules     *sync.Map
	databases   *sync.Map
	migrates    *migrate.EmbedDriver
	components  *sync.Map
	storage     storage.IStorage
}

// Init set config to application.
func Init(config config.IConfig) {
	GetInstance().Init(config)
}
func (app *App) Init(config config.IConfig) {
	app.config = config

	app.initLog(config.Sub("log"))
	app.initDefaultDatabase(config.Sub("default.database"))
	app.initDefaultStorage(config.Sub("default.storage"))
}

func (app *App) initDefaultStorage(config config.IConfig) {
	if config != nil {
		app.logger.Info("[application] init default storage...")
		var driver storage.IDriver
		switch config.GetString("driver") {
		case "local":
			opts := []func(*drivers.LocalDriver){}
			if config.GetString("access.scheme") != "" && config.GetString("access.domain") != "" {
				opts = append(opts, drivers.WithBaseUrl(fmt.Sprintf(
					"%s://%s",
					config.GetString("access.scheme"),
					config.GetString("access.domain"),
				)))
			}
			driver = drivers.NewLocal(storage.NewPath(config.GetString("path")), opts...)
		}

		if driver == nil {
			app.logger.Warn("[application] driver is not specified, disable storage")
			return
		}
		app.storage = storage.NewStorage(driver)
	}
}

func (app *App) initLog(config config.IConfig) {
	if config != nil {
		driver := config.GetString("driver", "zap")
		fmt.Printf("[application] init log using %s...\n", driver)
		switch driver {
		case "zap":
			app.logger = zap.NewZapLog()
		}
	} else {
		fmt.Printf("[application] log is not enabled\n")
	}
}

func (app *App) initDefaultDatabase(config config.IConfig) {
	if config != nil && config.GetBool("enable") {
		app.logger.Info("[application] init default database...")
		db, err := database.NewDatabase(strings.EncodeQuery(config.GetString("dsn")))
		if err != nil {
			panic(err)
		}
		app.RegisterDatabase("default", db)
	} else {
		app.logger.Warn("[application] default database is not enabled")
	}
}

func RegisterDatabase(name string, db interface{}) {
	GetInstance().RegisterDatabase(name, db)
}
func (app *App) RegisterDatabase(name string, db interface{}) {
	_, exists := app.databases.Load(name)
	if exists {
		panic(errx.New(fmt.Sprintf("database <%s> is already exists", name)))
	}
	app.databases.Store(name, db)
}

func (app *App) migrate() {
	if !app.config.GetBool("default.database.migrate") || !app.config.GetBool("default.database.enable") {
		app.logger.Info("[application] migrate is disabled")
		return
	}

	app.logger.Info("[application] migrating...")

	var (
		migrator migrate.IMigrator = &migrate.DefaultMigrator{}
		err      error
	)

	if err != nil {
		panic(err)
	}

	err = migrator.Migrate(strings.EncodeQuery(app.config.GetString("default.database.dsn")), app.migrates)
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
	app.modules.Range(func(key, module any) bool {
		app.logger.Info(fmt.Sprintf("boot module <%s>...", module.(IModule).Name()))
		err := module.(IModule).Boot(app)
		if err != nil {
			panic(err)
		}
		return true
	})
}

func (app *App) newDB(config map[string]interface{}) (db *gorm.DB, err error) {
	return database.NewDatabase(strings.EncodeQuery(config["dsn"].(string)))
}

// RegisterModule register module to application.
func RegisterModule(module interface{}) {
	GetInstance().RegisterModule(module)
}
func (app *App) RegisterModule(module interface{}) {
	m, ok := module.(IModule)
	if !ok {
		panic(errx.New("[application] invalid module"))
	}

	app.logger.Info(fmt.Sprintf("[application] register module <%s>...", m.Name()))
	err := m.Register(app)
	if err != nil {
		panic(errx.Wrap(err, "[application] register module failed"))
	}

	t := reflect.TypeOf(module)
	app.modules.Store(t, module)
}

// RegisterHttpServerRoute register module routes to application.
func RegisterHttpServerRoute(f func(serverEngine interface{}) error) {
	GetInstance().RegisterHttpServerRoute(f)
}
func (app *App) RegisterHttpServerRoute(f func(serverEngine interface{}) error) {
	svr, ok := app.httpServers.Load("default")
	if !ok {
		def := app.config.Sub("default.httpServer")
		if def != nil && def.GetBool("enable") {
			app.logger.Info("[application] create default http server...")
			svr = gin.NewGinServer(fmt.Sprintf(
				"%s:%d",
				def.GetString("host"),
				def.GetInt("port"),
			))
			app.httpServers.Store("default", svr)
		} else {
			app.logger.Info("[application] no default http server specifies")
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
		app.logger.Info(
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
		app.logger.Info(
			"[application] shutdown http server...",
			&log.Field{Name: "name", Value: key},
		)
		err := value.(server.IHTTPServer).Shutdown()
		if err != nil {
			app.logger.Error(
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

func Logger() log.ILog {
	return GetInstance().Logger()
}
func (app *App) Logger() log.ILog {
	return app.logger
}

func ModuleByName(nameOrType interface{}) interface{} {
	return GetInstance().Module(nameOrType)
}
func Module[T IModule]() T {
	var typePtr *T
	t := reflect.TypeOf(typePtr).Elem()
	return GetInstance().Module(t).(T)
}
func (app *App) Module(nameOrType interface{}) (module interface{}) {
	switch key := nameOrType.(type) {
	case string:
		app.modules.Range(func(k, value any) bool {
			if value.(IModule).Name() == key {
				module = value
				return false
			}
			return true
		})
	case reflect.Type:
		module, _ = app.modules.Load(key)
	default:
		t := reflect.TypeOf(nameOrType)
		module, _ = app.modules.Load(t)
	}

	return
}

func Config() config.IConfig {
	return GetInstance().Config()
}
func (app *App) Config() config.IConfig {
	return app.config
}

func SetComponent(key any, value any) {
	GetInstance().SetComponent(key, value)
}
func (app *App) SetComponent(key any, value any) {
	app.components.Store(key, value)
}

func GetComponent(key any) any {
	return GetInstance().GetComponent(key)
}
func (app *App) GetComponent(key any) any {
	v, _ := app.components.Load(key)
	return v
}

func Storage() storage.IStorage {
	return GetInstance().Storage()
}
func (app *App) Storage() storage.IStorage {
	return app.storage
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
