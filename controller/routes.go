package controller

import (
	"github.com/labstack/echo"
)

func Routing(e *echo.Echo) {
	g := e.Group("/task")
	g.GET("/favorite_tweets", favoriteTweets)
}
