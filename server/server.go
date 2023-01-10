package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/erfanshekari/go-talk/config"

	"github.com/erfanshekari/go-talk/routes"
	"github.com/erfanshekari/go-talk/test"

	"github.com/erfanshekari/go-talk/internal/global"
	ujwt "github.com/erfanshekari/go-talk/utils/jwt"
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

func (s *Server) registerAuthMiddlewares(g *echo.Group) {
	// Adding jwt auth middleware
	g.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(global.GetInstance(nil).SecretKey),
		ContextKey: "user",
	}))

	// Change jwt.Token to models.User
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get("user").(*jwt.Token)
			user, err := ujwt.GetUserFromJWT(token)
			if err != nil {
				return c.JSON(http.StatusForbidden, nil)
			}
			c.Set("user", user)
			return next(c)
		}
	})
}

func (s *Server) Listen() {

	e := echo.New() // init server

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
	s.registerAuthMiddlewares(rest)

	// adding validator
	e.Validator = GetValidator()

	// register routes
	s.registerRoutes(e)

	// register jwt protected routes
	s.registerProtectedRoutes(rest)

	// starting server
	addr := s.Config.Server.Host + ":" + strconv.FormatInt(int64(s.Config.Server.Port), 10)
	err := e.Start(addr)
	if err != nil {
		log.Println(err)
	}
}
