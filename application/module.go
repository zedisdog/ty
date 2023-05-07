package application

import (
	"fmt"
	"github.com/zedisdog/ty/errx"
	"reflect"
)

type IHasModule interface {
	Module(nameOrType interface{}) (module interface{})
	RegisterModule(module IModule)
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

type IModule interface {
	Name() string
	// Register registers resource to application. e.g: route used by default http server
	Register() error
	// Boot starts module's own sub process.
	Boot() error
}
