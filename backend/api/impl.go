package api

import (
	"bar/autogen"
	"bar/internal/db"
	"os"
	"time"

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

func (s *Server) Serve() error {
	e := echo.New()

	// Logger
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.INFO)
	log.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	e.Logger = log.Logger()
	e.Use(middleware.Logger())

	autogen.RegisterHandlers(e, s)

	if err := e.Start(":8080"); err != nil {
		return err
	}

	return nil
}
