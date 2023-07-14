package api

import (
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/db"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	middleware "github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

type Server struct {
	db.DBackend
}

func NewServer(db db.DBackend) *Server {
	s := &Server{
		db,
	}
	return s
}

func (s *Server) Serve(c *config.Config) error {
	e := echo.New()

	// Logger
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(middleware.Logger())

	userStore := sessions.NewCookieStore([]byte(c.ApiConfig.SessionSecret))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer context.Clear(c.Request())
			c.Set("userStore", userStore)
			return next(c)
		}
	})

	adminStore := sessions.NewCookieStore([]byte(c.ApiConfig.AdminSessionSecret))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer context.Clear(c.Request())
			c.Set("adminStore", adminStore)
			return next(c)
		}
	})

	// You can use h for intellisense and get the handlers' names
	h := autogen.NewStrictHandler(s, []autogen.StrictMiddlewareFunc{
		s.AuthMiddleware,
	})

	autogen.RegisterHandlers(e, h)

	if err := e.Start(c.ApiConfig.Port); err != nil {
		return err
	}

	return nil
}
