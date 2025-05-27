package api

import (
	"net/http"

	"github.com/jordyvanvorselen/go-templ-htmx-vercel-template/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/users", controller.UserHandler{}.Index)

	e.ServeHTTP(w, r)
}
