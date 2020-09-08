package main

import (
	"net/http"
	"os"

	"github.com/fabienbellanger/goutils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type server struct {
	store  store
	router *echo.Echo
	mode   string
}

func newServer() *server {
	s := &server{
		router: echo.New(),
		mode:   "production",
	}

	s.initHTTPServer()
	s.routes()
	s.initPprof()
	s.errorHandling()

	return s
}

func (s *server) initHTTPServer() {
	// Mode
	// ----
	s.mode = viper.GetString("environment")

	// Recover
	// -------
	s.router.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisableStackAll:   true,
		DisablePrintStack: true,
		Skipper:           middleware.DefaultSkipper,
	}))

	// Logger
	// ------
	if s.mode != "production" {
		s.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format:           "${time_custom} | ${remote_ip}\t| ${status} | ${latency_human}\t| ${method}\t| ${uri}\n",
			CustomTimeFormat: "2006-01-02 15:04:05",
			Output:           os.Stderr,
		}))
	}

	// Security
	// --------
	s.router.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         3600,
		// ContentSecurityPolicy: "default-src 'self'",
	}))

	// CORS
	// ----
	s.router.Use(middleware.CORS())
}

func (s *server) initPprof() {
	basicAuthUsername := viper.GetString("debug.basicAuthUsername")
	basicAuthPassword := viper.GetString("debug.basicAuthPassword")

	if viper.GetBool("debug.pprof") && basicAuthUsername != "" && basicAuthPassword != "" {
		g := s.router.Group("private")

		// Protection des routes par une Basic Auth
		// ----------------------------------------
		g.Use(middleware.BasicAuth(
			func(username, password string, c echo.Context) (bool, error) {
				if username == basicAuthUsername && password == basicAuthPassword {
					return true, nil
				}
				return false, nil
			},
		))

		s.privatePprofRoutes(g)
	}
}

func (s *server) errorHandling() {
	s.router.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if httpError, ok := err.(*echo.HTTPError); ok {
			code = httpError.Code
		}

		msg := goutils.GetHTTPStatusMessage(code)
		if msg == "" {
			msg = "Unknown Error"
		}

		c.JSON(code, map[string]string{"message": msg})
	}
}
