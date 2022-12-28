package context

import (
	"github.com/erfanshekari/go-talk/models"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	User *models.User `json:"user"`
}
