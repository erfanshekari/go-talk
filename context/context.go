package context

import (
	"github.com/erfanshekari/go-talk/config"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	ServerConfig *config.Config
}
