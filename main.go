package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"html/template"
	"path/filepath"
	"net/http"
)


type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, struct {
		Data    interface{}
		Context echo.Context
	}{data, c})

	if err != nil {
		panic(err)
	}

	return nil
}


func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.New("templater").Funcs(nil).ParseGlob(filepath.Join("templates/", "*.html"))),
	}

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "found")
	})

	panic(e.Start(":8080"))
}
