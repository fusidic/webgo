package main

import (
	"net/http"

	"github.com/fusidic/webgo/pkg/base"
)

func main() {
	r := base.New()
	r.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>hi there</h1>")
	})
	r.GET("/hello", func(c *base.Context) {
		c.String(http.StatusOK, "hello %s, Path: %s, Method: %s\n", c.Query("name"), c.Path, c.Method)
	})

	r.POST("/login", func(c *base.Context) {
		c.JSON(http.StatusOK, base.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":9999")
}
