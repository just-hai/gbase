package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/just-hai/gbase/pkg/config"
	"github.com/just-hai/gbase/pkg/server/handler"
	mware "github.com/just-hai/gbase/pkg/server/middleware"
	"github.com/just-hai/gbase/pkg/server/router"
	"github.com/just-hai/gbase/pkg/swagger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
	*router.RouterGroup
	conf *config.Config
}

func New() *Server {
	return NewWithConfig(&config.Conf)
}

func NewWithConfig(conf *config.Config) *Server {
	ech := echo.New()
	server := &Server{
		Echo:        ech,
		RouterGroup: &router.RouterGroup{EchoGroup: ech.Group("")},
		conf:        conf,
	}
	return server
}

func (server *Server) Start() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server.Echo.Use(middleware.Recover())
	server.Echo.HTTPErrorHandler = handler.HTTPErrorHandler
	server.Echo.Use(mware.RequestLogger())

	if config.Conf.Swagger {
		server.Echo.GET("/swagger.json", func(c echo.Context) error {
			swaggerJson, _ := router.GetSwaggerJSON()
			c.Response().Write([]byte(swaggerJson))
			return nil
		})
		server.Echo.StaticFS("/swagger", echo.MustSubFS(swagger.UI, "ui"))
	}
	go func() {
		err := server.Echo.Start(server.conf.Port)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Echo.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
