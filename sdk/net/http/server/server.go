package server

type IHTTPServer[T any] interface {
	// RegisterRoutes register routes to server
	RegisterRoutes(func(engine T) error) error
	// Run starts server with non-blocking.
	Run()
	// Shutdown shutdowns server. Graceful will be best.
	Shutdown() error
}
