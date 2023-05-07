package application

type IHasComponent interface {
	RegisterComponent(name any, value any)
	Component(name any) any
}

func RegisterComponent(name any, value any) {
	GetInstance().RegisterComponent(name, value)
}
func (app *App) RegisterComponent(name any, value any) {
	app.components.Store(name, value)
}

func Component[T any](name any) T {
	return GetInstance().Component(name).(T)
}
func (app *App) Component(name any) any {
	v, _ := app.components.Load(name)
	return v
}
