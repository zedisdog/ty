package application

import (
	"errors"
	"fmt"
	"io/fs"

	migrate2 "github.com/golang-migrate/migrate/v4"
	"github.com/zedisdog/ty/database/migrate"
	"github.com/zedisdog/ty/strings"
)

type IHasMigrator interface {
	RegisterMigrator(name string, migrator migrate.IMigrator)
	Migrator(name ...string) migrate.IMigrator
	RegisterMigration(fs fs.FS)
}

func RegisterMigrator(name string, migrator migrate.IMigrator) {
	GetInstance().RegisterMigrator(name, migrator)
}
func (app *App) RegisterMigrator(name string, migrator migrate.IMigrator) {
	app.migrators.Store(name, migrator)
}

func Migrator(name ...string) migrate.IMigrator {
	return GetInstance().Migrator(name...)
}
func (app *App) Migrator(name ...string) migrate.IMigrator {
	return app.getValueOrDefault(app.migrators, name...).(migrate.IMigrator)
}

func RegisterMigration(fs fs.FS) {
	GetInstance().RegisterMigration(fs)
}
func (app *App) RegisterMigration(fs fs.FS) {
	app.Migrator().GetSourceInstance().(*migrate.EmbedDriver).Add(fs)
}

func (app *App) initDefaultMigrator() {
	if !app.config.GetBool("default.database.migrate") {
		return
	}

	migrator := &migrate.DefaultMigrator{}

	migrator.SetDatabaseURL(strings.EncodeQuery(app.config.GetString("default.database.dsn")))

	migrator.SetSourceInstance(migrate.NewFsDriver())

	app.RegisterMigrator("default", migrator)
}

func (app *App) migrate() {
	app.migrators.Range(func(name, value any) bool {
		app.logger.Info(fmt.Sprintf("migrator %s up...", name.(string)))
		err := value.(migrate.IMigrator).Migrate()
		if err != nil && !errors.Is(err, migrate2.ErrNoChange) {
			panic(err)
		}
		return true
	})
}
