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
	Config *config.Config
}

func (s *Server) registerRoutes(e *echo.Echo) {
	for _, route := range routes.Routes {
		route(e)
	}
}

func (s *Server) registerProtectedRoutes(g *echo.Group) {
	for _, route := range routes.ProtectedRoutes {
		route(g)
	}
}

func (s *Server) Listen() {

	e := echo.New() // init server

	// Change Default Context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ctx.Context{
				Context:      c,
				ServerConfig: s.Config,
			}
			return next(cc)
		}
	})

	// register echo logger and recover middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if s.Config.Server.Debug {
		// performing build tests if necessary.
		test.RegisterTest(e, s.Config.Server.LazyDebug)

		err := godotenv.Load()

		if err != nil {
			log.Fatal(err)
		}
	}

	// Adding Throttle Middleware
	e.Use(middleware.RateLimiterWithConfig(ThrottleConfig))

	// adding cors config
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	rest := e.Group("/rest")

	// Adding jwt auth middleware
	secretKey := os.Getenv("SECRET_KEY")
	rest.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secretKey),
		ContextKey: "user",
	}))

	// Register user to context
	rest.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*ctx.Context)
			token := cc.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			cc.Set("user", &models.User{
				UserID: strconv.FormatInt(
					int64(int((claims["user_id"]).(float64))),
					10,
				)})
			return next(cc)
		}
	})

	// adding validator
	e.Validator = GetValidator()

	// register routes
	s.registerRoutes(e)

	// register jwt protected routes
	s.registerProtectedRoutes(rest)

	// starting server
	addr := s.Config.Server.Host + ":" + strconv.FormatInt(int64(s.Config.Server.Port), 10)
	e.Start(addr)
}
