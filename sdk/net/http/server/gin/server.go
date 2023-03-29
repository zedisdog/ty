package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/sdk/net/http/server"
	"net/http"
	"time"
)

func NewGinServer(addr string) server.IHTTPServer {
	return &Server{
		engine: gin.Default(),
		srv: &http.Server{
			Addr: addr,
		},
	}
}

type Server struct {
	engine *gin.Engine
	srv    *http.Server
	logger log.ILog
}

func (s Server) RegisterRoutes(f func(serverEngine interface{}) error) error {
	return f(s.engine)
}

func (s Server) Run() {
	s.srv.Handler = s.engine
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("server shutdown", log.NewField("err", err))
		}
	}()
}

func (s Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

func (s Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.engine.ServeHTTP(res, req)
}
