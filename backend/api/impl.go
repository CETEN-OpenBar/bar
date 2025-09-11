package api

import (
	"bar/autogen"
	"bar/autogen/helloasso"
	"bar/internal/config"
	"bar/internal/db"
	"context"
	"os"
	"time"

	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	middlewareL "github.com/neko-neko/echo-logrus/v2"
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/sirupsen/logrus"
)

type Server struct {
	db.DBackend
	HelloAssoClient *helloasso.ClientWithResponses
}

func NewHelloAssoClient() (*helloasso.ClientWithResponses, error) {
	c := config.GetConfig()

	client, err := helloasso.NewClient(c.HelloAssoConfig.URL + "/v5")
	if err != nil {
		return nil, err
	}

	_, err = client.GetToken()
	if err != nil {
		return nil, err
	}

	return &helloasso.ClientWithResponses{ClientInterface: client}, nil
}

// Try to setup the helloasso client for the given server
func (s *Server) SetupHelloAsso() error {
	client, err := NewHelloAssoClient()
	if err != nil || client == nil {
		logrus.Error("Error initializing HelloAsso client : ", err);
		return err;
	}
	s.HelloAssoClient = client;
	logrus.Info("HelloAsso client initialized successfully !")
	return nil;
}

func NewServer(db db.DBackend) *Server {

	s := &Server{
		db,
		nil,
	}

	err := s.SetupHelloAsso();
	if err != nil {
		// Retry HelloAsso setup later and run without remote refills
		logrus.Info("Could not initialize HelloAsso client, trying again every 2 minutes")
		go func() {
			// API rate limits for authentication endpoints are 
    		// 10 calls per 10 seconds
    		// 20 calls per 10 minutes
    		// 50 calls per hour
			ticker := time.NewTicker(2 * time.Minute)
			defer ticker.Stop()
			
			for range ticker.C {
				err := s.SetupHelloAsso()
				if err == nil {
					// Setup was successful, stop the task
					return;
				}
			}
		}()
	}

	return s
}

func (s *Server) getUserSess(c echo.Context) *sessions.Session {
	return c.Get("userSess").(*sessions.Session)
}

func (s *Server) getAdminSess(c echo.Context) *sessions.Session {
	return c.Get("adminSess").(*sessions.Session)
}
func (s *Server) getOnboardSess(c echo.Context) *sessions.Session {
	return c.Get("onBoardSess").(*sessions.Session)
}

// The remote refill subsystem can only be used if the helloasso client is available
func (s *Server) remoteRefillsAvailable() bool {
	return s.HelloAssoClient != nil;
}

func (s *Server) Serve(c *config.Config) error {
	e := echo.New()

	// Logger
	log.Logger().SetOutput(os.Stdout)
	log.Logger().SetLevel(echoLog.WARN)
	e.Logger = log.Logger()
	e.Use(middlewareL.Logger())
	e.Use(middleware.BodyLimit("15M"))

	userStore := sessions.NewCookieStore([]byte(c.ApiConfig.SessionSecret))
	adminStore := sessions.NewCookieStore([]byte(c.ApiConfig.AdminSessionSecret))
	onBoardStore := sessions.NewCookieStore([]byte(c.ApiConfig.AdminSessionSecret))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer gorillaContext.Clear(c.Request())
			c.Set("userStore", userStore)
			c.Set("adminStore", adminStore)
			c.Set("onBoardStore", onBoardStore)
			return next(c)
		}
	})

	// CORS
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", c.Request().Header.Get("Origin"))
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Cookie, Cookies, X-Local-Token")
			return next(c)
		}
	})

	// Fake timeout
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		time.Sleep(1000 * time.Millisecond)
	// 		return next(c)
	// 	}
	// })

	// You can use h for intellisense and get the handlers' names
	e.Use(s.AuthMiddleware)

	autogen.RegisterHandlers(e, s)

	// Start the HelloAsso processing runner
	go func() {
		ticker := time.NewTicker(c.HelloAssoConfig.CheckoutProcessingInterval)
		defer ticker.Stop()
		
		ctx := context.Background()

		for range ticker.C {
			if s.remoteRefillsAvailable() {
				s.ProcessStartedRefills(ctx)
			}	
		}
	}()

	if err := e.Start(c.ApiConfig.Port); err != nil {
		return err
	}

	return nil
}
