package application

type IModule interface {
	Name() string
	// Register registers resource to application. e.g: route used by default http server
	Register(application IApplication) error
	// Boot starts module's own sub process.
	Boot(application IApplication) error
}
