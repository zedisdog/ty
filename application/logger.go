package application

import (
	"fmt"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/log/zap"
)

type IHasLogger interface {
	Logger() log.ILog
}

func (app *App) initLog() {
	config := app.config.Sub("log")
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

func Logger() log.ILog {
	return GetInstance().Logger()
}
func (app *App) Logger() log.ILog {
	return app.logger
}
