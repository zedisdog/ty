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

func Component[T any](name ...any) T {
	if len(name) > 0 {
		return GetInstance().Component(name[0]).(T)
	} else {
		var found T
		GetInstance().components.Range(func(key, value any) bool {
			if f, ok := value.(T); ok {
				found = f
				return false
			}

			return true
		})
		return found
	}
}
func (app *App) Component(name any) any {
	v, _ := app.components.Load(name)
	return v
}
