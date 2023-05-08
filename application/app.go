package application

import (
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/log"
	"net/http"
	"sync"
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
			onStop:      make([]func(), 0),
			seeders:     make([]func() error, 0),
			storages:    new(sync.Map),
			migrators:   new(sync.Map),
		}
	})

	return instance
}

type IApplication interface {
	ILifetime
	IHasComponent
	ICanTest
	IHasScheduler
	IHasDatabase
	IHasModule
	IHasConfig
	IHasStorage
	IHasLogger
	IHasHttpServer
	IHasMigrator
}

var _ IApplication = (*App)(nil)

type App struct {
	config      *config.Config
	httpServers *sync.Map
	logger      log.ILog
	modules     *sync.Map
	databases   *sync.Map
	migrators   *sync.Map
	seeders     []func() error
	components  *sync.Map
	storages    *sync.Map
	onStop      []func()

	/*************test**************/

	header http.Header
}

func (app *App) getValueOrDefault(m *sync.Map, name ...string) (value any) {
	if len(name) > 0 {
		value, _ = m.Load(name[0])
	}
	if value != nil {
		return
	}

	value, ok := m.Load("default")
	if ok {
		return
	}

	count := 0
	app.databases.Range(func(key, v any) bool {
		count++
		value = v
		return true
	})

	if count == 1 {
		return
	} else {
		return nil
	}
}
