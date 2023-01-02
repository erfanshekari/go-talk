package server

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
)

var DefaultCORSConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
	},
}
