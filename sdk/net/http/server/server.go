package server

type IHTTPServer interface {
	// RegisterRoutes register routes to server
	RegisterRoutes(func(serverEngine interface{}) error) error
	// Run starts server with non-blocking.
	Run()
	// Shutdown shutdowns server. Graceful will be best.
	Shutdown() error
}
