package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"html/template"
	"path/filepath"
	"net/http"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"./models"
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

var db *sqlx.DB

func connectToDB() {
	var err error

	db, err = sqlx.Open("sqlite3", "beosztas2.db")
	if err != nil {
		panic(err)
	}
}

type paginationData struct {
	_page      int
	_perPage   int
	_sortDir   string
	_sortField string
}

func parsePaginationData(c echo.Context) (*paginationData, error) {
	res := &paginationData{}
	var err error

	_page := c.QueryParam("_page")
	_perPage := c.QueryParam("_perPage")

	res._sortDir = c.QueryParam("_sortDir")
	res._sortField = c.QueryParam("_sortField")

	res._page, err = strconv.Atoi(_page)
	if err != nil {
		return nil, err
	}

	res._perPage, err = strconv.Atoi(_perPage)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	connectToDB()
	models.SetDB(db)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.New("templater").Funcs(nil).ParseGlob(filepath.Join("templates/", "*.html"))),
	}

	e.Renderer = t

	e.Static("/static", "public")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	})

	admin := e.Group("/admin", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "admin" {
			return true, nil
		}
		return false, nil
	}))

	admin.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "admin.html", nil)
	})

	admin.GET("/classes", getClasses)
	admin.GET("/classes/:id", getClass)
	admin.PUT("/classes/:id", putClass)
	admin.POST("/classes", postClass)
	admin.DELETE("/classes/:id", deleteClass)

	panic(e.Start(":8080"))
}
