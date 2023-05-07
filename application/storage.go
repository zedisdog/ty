package application

import (
	"fmt"
	"github.com/zedisdog/ty/storage"
	"github.com/zedisdog/ty/storage/drivers"
)

type IHasStorage interface {
	Storage(name ...string) any
	RegisterStorage(key string, store any)
}

func RegisterStorage(key string, store any) {
	GetInstance().RegisterStorage(key, store)
}
func (app *App) RegisterStorage(key string, store any) {
	app.storages.Store(key, store)
}

// Storage return storage instance by name.
//
// if no name, return default instance.
//
// if no default instance, and there only one instance, return it.
//
// if there are more than one instance, return nil, cause can not determine which instance should be returned.
func Storage[T any](name ...string) T {
	return GetInstance().Storage(name...).(T)
}
func (app *App) Storage(name ...string) any {
	return app.getValueOrDefault(app.storages, name...)
}

func (app *App) initDefaultStorage() {
	config := app.config.Sub("default.storage")
	if config != nil {
		app.logger.Info("[application] init default storage...")
		var driver storage.IDriver
		switch config.GetString("driver") {
		case "local":
			var opts []func(*drivers.LocalDriver)
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
		app.RegisterStorage("default", storage.NewStorage(driver))
	}
}
