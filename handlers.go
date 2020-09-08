package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *server) handlerHome(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World 👋!")
}

func (s *server) handlerBigJSON(c echo.Context) error {
	type User struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"-"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
	}
	var users []User
	for i := 0; i < 100000; i++ {
		users = append(users, User{
			ID:        i + 1,
			Username:  "My Username",
			Lastname:  "My Lastname",
			Firstname: "My Firstname",
		})
	}
	return c.JSON(http.StatusOK, &users)
}

func (s *server) handlerBigJSONStream(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
