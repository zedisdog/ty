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
	"syscall"
)

type ILifetime interface {
	Init(config config.IConfig)
	RegisterDatabase(name string, db interface{})
	RegisterModule(module interface{})
	RegisterHttpServerRoute(f func(serverEngine interface{}) error)
	RegisterMigrate(fs *embed.FS)
	RegisterSeeder(seeders ...func(app IApplication) error)
	RegisterStopFunc(f func())
	Boot()
	Run()
	Wait(closeFunc ...func())
}

/*********************************init*********************************************/

// Init set config to application.
func Init(config config.IConfig) {
	GetInstance().Init(config)
}
func (app *App) Init(config config.IConfig) {
	app.config = config

	app.initLog(config.Sub("log"))
	app.initDefaultDatabase(config.Sub("default.database"))
	app.initDefaultStorage(config.Sub("default.storage"))
	app.RegisterStopFunc(app.stop)
}

func (app *App) initDefaultStorage(config config.IConfig) {
	if config != nil && config.GetBool("enable") {
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

/****************************************register***********************************/

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

// RegisterModule register module to application.
func RegisterModule(module IModule) {
	GetInstance().RegisterModule(module)
}
func (app *App) RegisterModule(module IModule) {

	app.logger.Info(fmt.Sprintf("[application] register module <%s>...", module.Name()))
	err := module.Register()
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
			svr = gin.NewGinServer(
				fmt.Sprintf(
					"%s:%d",
					def.GetString("host"),
					def.GetInt("port"),
				),
				gin.EnablePprof(def.GetBool("enablePprof")),
			)
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

func RegisterMigrate(fs *embed.FS) {
	GetInstance().migrates.Add(fs)
}
func (app *App) RegisterMigrate(fs *embed.FS) {
	app.migrates.Add(fs)
}

func RegisterSeeder(seeders ...func() error) {
	GetInstance().RegisterSeeder(seeders...)
}
func (app *App) RegisterSeeder(seeders ...func() error) {
	if len(seeders) < 1 {
		return
	}

	app.seeders = append(app.seeders, seeders...)
}

func RegisterStopFunc(f func()) {
	GetInstance().RegisterStopFunc(f)
}
func (app *App) RegisterStopFunc(f func()) {
	app.onStop = append(app.onStop, f)
}

/******************************************run**********************************************/

// Boot boots the application.
func Boot() {
	GetInstance().Boot()
}
func (app *App) Boot() {
	app.migrate()
	app.seed()
	app.modules.Range(func(key, module any) bool {
		app.logger.Info(fmt.Sprintf("boot module <%s>...", module.(IModule).Name()))
		err := module.(IModule).Boot()
		if err != nil {
			panic(err)
		}
		return true
	})
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

func (app *App) stop() {
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
	app.CloseScheduler()
}

func Wait(closeFunc ...func()) {
	GetInstance().Wait(closeFunc...)
}
func (app *App) Wait(closeFunc ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	for _, cls := range append(app.onStop, closeFunc...) {
		cls()
	}
}

/******************************************************************************************/

func (app *App) migrate() {
	if !app.config.GetBool("default.database.migrate") || !app.config.GetBool("default.database.enable") {
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

func (app *App) seed() {
	if !app.config.GetBool("default.database.migrate") || !app.config.GetBool("default.database.enable") {
		return
	}

	app.logger.Info("[application] seeding...")

	for _, seeder := range app.seeders {
		err := seeder()
		if err != nil {
			panic(err)
		}
	}
}

func (app *App) newDB(config map[string]interface{}) (db *gorm.DB, err error) {
	return database.NewDatabase(strings.EncodeQuery(config["dsn"].(string)))
}
