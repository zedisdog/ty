package application

import (
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/storage"
	"reflect"
)

type IHasDatabase interface {
	Database(name string) (db interface{})
}

type IHasComponent interface {
	SetComponent(key any, value any)
	GetComponent(key any) any
	Logger() log.ILog
	Module(nameOrType interface{}) (module interface{})
	Config() config.IConfig
	Storage() storage.IStorage
	IHasDatabase
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

func Storage() storage.IStorage {
	return GetInstance().Storage()
}
func (app *App) Storage() storage.IStorage {
	return app.storage
}
