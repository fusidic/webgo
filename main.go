package main

import (
	"net/http"

	"github.com/fusidic/webgo/pkg/base"
)

func main() {
	e := base.New()
	e.GET("/index", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>Index</h1>")
	})

	v1 := e.Group("/v1")

	v1.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "<h1>hi in v1</h1>")
	})
	v1.GET("/hello", func(c *base.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	v2 := e.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *base.Context) {
			c.String(http.StatusOK, "hello %s, your're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *base.Context) {
			c.JSON(http.StatusOK, base.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	e.Run(":9999")
}
