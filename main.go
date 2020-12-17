package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/fusidic/webgo/pkg/base"
)

type student struct {
	Name string
	Age  int8
}

// V2middle ...
func V2middle() base.HandlerFunc {
	return func(c *base.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// FormatAsDate format time
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%-2d", year, month, day)
}

func main() {
	e := base.New()
	// global middleware
	e.Use(base.Logger(), base.Recovery())

	// 用户定义渲染函数
	e.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 其他默认渲染
	e.LoadHTMLGlob("templates/*")

	// 加入静态地址与对应的handler
	e.Static("/assets", "./static")

	stu1 := &student{Name: "fusidic", Age: 20}
	stu2 := &student{Name: "arithbar", Age: 21}
	// HTML中加载所有e.htmlTemplates
	e.GET("/", func(c *base.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	e.GET("/students", func(c *base.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", base.H{
			"title":  "base",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	e.GET("/date", func(c *base.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", base.H{
			"title": "base",
			"now":   time.Date(2020, 12, 20, 0, 0, 0, 0, time.UTC),
		})
	})

	e.Static("/assets", "/Users/fusidic/Documents/github.com/fusidic/webgo/static")

	// group v1
	v1 := e.Group("/v1")

	v1.GET("/hello", func(c *base.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	// group v2
	v2 := e.Group("/v2")
	// apply middleware V2middle to v2
	v2.Use(V2middle())
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

	// index out of range for testing Recovery()
	e.GET("/panic", func(c *base.Context) {
		names := []string{"fusidic"}
		c.String(http.StatusOK, names[100])
	})

	e.Run(":9999")
}
