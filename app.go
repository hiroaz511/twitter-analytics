package main

import (
	"github.com/labstack/echo"
	"google.golang.org/appengine"
	"net/http"
	"twitter-analytics/controller"
)

func main() {
	e := echo.New()
	controller.Routing(e)

	http.Handle("/", e)

	appengine.Main()
}

