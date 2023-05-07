package application

import (
	"fmt"
	"github.com/zedisdog/ty/database"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/strings"
)

type IHasDatabase interface {
	Database(name ...string) (db interface{})
	RegisterDatabase(name string, db interface{})
}

func (app *App) initDefaultDatabase() {
	config := app.config.Sub("default.database")
	if config != nil {
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

// Database return database instance by name.
//
// if no name, return default instance.
//
// if no default instance, and there only one instance, return it.
//
// if there are more than one instance, return nil, cause can not determine which instance should be returned.
func Database[T any](name ...string) T {
	return GetInstance().Database(name...).(T)
}
func (app *App) Database(name ...string) (db interface{}) {
	return app.getValueOrDefault(app.databases, name...)
}
