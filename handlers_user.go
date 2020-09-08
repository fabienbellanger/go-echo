package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) handlerGetUser(c echo.Context) error {
	u := s.store.user().getUser()

	return c.JSON(http.StatusOK, &u)
}
