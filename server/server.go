package server

import (
	"log"
	"os"
	"strconv"

	"github.com/erfanshekari/go-talk/config"
	ctx "github.com/erfanshekari/go-talk/context"
	"github.com/erfanshekari/go-talk/models"
	"github.com/erfanshekari/go-talk/routes"
	"github.com/erfanshekari/go-talk/test"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config *config.ConfigAtrs
}

func (s *Server) registerRoutes(e *echo.Echo) {
	for _, route := range routes.Routes {
		route(e)
	}
}

func (s *Server) Listen() {

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ctx.Context{Context: c, User: nil}
			return next(cc)
		}
	})

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
		ContextKey: "user",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*ctx.Context)
			token := cc.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			cc.User = &models.User{
				UserID: strconv.FormatInt(
					int64(int((claims["user_id"]).(float64))),
					10,
				)}
			return next(cc)
		}
	})

	s.registerRoutes(e)

	addr := s.Config.IP + ":" + strconv.FormatInt(int64(s.Config.Port), 10)
	e.Start(addr)
}
