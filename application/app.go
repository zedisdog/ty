package application

import (
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/database/migrate"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/storage"
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
			migrates:    migrate.NewFsDriver(),
			onStop:      make([]func(), 0),
			seeders:     make([]func(app IApplication) error, 0),
		}
	})

	return instance
}

type IApplication interface {
	ILifetime
	IHasComponent
	ICanTest
	IHasScheduler
}

type App struct {
	config      config.IConfig
	httpServers *sync.Map
	logger      log.ILog
	modules     *sync.Map
	databases   *sync.Map
	migrates    *migrate.EmbedDriver
	seeders     []func(app IApplication) error
	components  *sync.Map
	storage     storage.IStorage
	onStop      []func()

	/*************test**************/

	header http.Header
}
