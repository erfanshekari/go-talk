package server

import (
	"log"
	"os"
	"strconv"

	"github.com/erfanshekari/go-talk/config"
	"github.com/erfanshekari/go-talk/server/routes"
	"github.com/erfanshekari/go-talk/test"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config *config.ConfigAtrs
}

func (s *Server) registerRoutes(e *echo.Echo) {
	routes.RegisterRoutes(e)
}

func (s *Server) Listen() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if s.Config.Debug {
		test.RegisterTest(e, s.Config.DebugLazy)
	}

	if s.Config.Debug {
		err := godotenv.Load()

		if err != nil {
			log.Fatal(err)
		}
	}

	secretKey := os.Getenv("SECRET_KEY")

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secretKey),
	}))

	s.registerRoutes(e)

	addr := s.Config.IP + ":" + strconv.FormatInt(int64(s.Config.Port), 10)
	e.Start(addr)
}
