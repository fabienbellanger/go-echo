package main

import (
	"encoding/json"
	"io"
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
	type User struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"-"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
	}

	response := c.Response()
	response.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	response.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(response, "{"); err != nil {
		return err
	}

	encoder := json.NewEncoder(response)
	for i := 0; i < 100000; i++ {
		if i > 0 {
			if _, err := io.WriteString(response, ","); err != nil {
				return err
			}
		}

		user := User{
			ID:        i + 1,
			Username:  "My Username",
			Lastname:  "My Lastname",
			Firstname: "My Firstname",
		}

		if err := encoder.Encode(user); err != nil {
			return err
		}

		i++
	}

	if _, err := io.WriteString(response, "}"); err != nil {
		return err
	}

	return nil
}
