package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/sdk/net/http/middlewares"
	"github.com/zedisdog/ty/sdk/net/http/server"
	"net/http"
	"net/http/pprof"
	"time"
)

func EnablePprof(enable bool) func(svr *Server) {
	return func(svr *Server) {
		svr.enablePprof = enable
	}
}

func NewGinServer(addr string, options ...func(svr *Server)) server.IHTTPServer {
	svr := &Server{
		engine: gin.Default(),
		srv: &http.Server{
			Addr: addr,
		},
		enablePprof: false,
	}

	svr.engine.Use(middlewares.Cros)

	for _, option := range options {
		option(svr)
	}

	if svr.enablePprof {
		svr.engine.GET("/debug/pprof/*action", func(c *gin.Context) {
			switch c.Param("action") {
			case "/cmdline":
				pprof.Cmdline(c.Writer, c.Request)
			case "/profile":
				pprof.Profile(c.Writer, c.Request)
			case "/symbol":
				pprof.Symbol(c.Writer, c.Request)
			case "/trace":
				pprof.Trace(c.Writer, c.Request)
			default:
				pprof.Index(c.Writer, c.Request)
			}
		})
	}

	return svr
}

type Server struct {
	engine      *gin.Engine
	srv         *http.Server
	logger      log.ILog
	enablePprof bool
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
