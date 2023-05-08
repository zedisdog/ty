package application

type IHasSeeders interface {
	RegisterSeeder(seeders ...func() error)
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

func (app *App) seed() {
	app.logger.Info("[application] seeding...")

	for _, seeder := range app.seeders {
		err := seeder()
		if err != nil {
			panic(err)
		}
	}
}
