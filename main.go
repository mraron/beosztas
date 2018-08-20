package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"html/template"
	"path/filepath"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"./models"
	"github.com/jinzhu/gorm"
	"fmt"
	"encoding/json"
	"time"
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

var db *gorm.DB

func connectToDB() {
	var err error

	db, err = gorm.Open("sqlite3", "beosztas.db")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	db.AutoMigrate(&models.Class{})
	db.AutoMigrate(&models.Student{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Place{})
	db.AutoMigrate(&models.Participation{})
}

type paginationData struct {
	_filters map[string]string
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
		res._page = 1
		//return nil, err
	}

	res._perPage, err = strconv.Atoi(_perPage)
	if err != nil {
		res._perPage = 1000000
	}

	_filters := c.QueryParam("_filters")
	res._filters = make(map[string]string)

	json.Unmarshal([]byte(_filters), &res._filters)

	return res, nil
}


func parseUint(str string) (uint, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

func main() {
	connectToDB()
	models.SetDB(db)

	st := new(models.Student)
	db.First(st)
	fmt.Println(st)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.New("templater").Funcs(template.FuncMap{
			"localtime": func(t time.Time) string {
				loc, err := time.LoadLocation("Europe/Budapest")
				if err != nil {
					panic(err)
				}

				return t.In(loc).Format("2006.Jan.2 15:04:05")
			},
		}).ParseGlob(filepath.Join("templates/", "*.html"))),
	}

	e.Renderer = t

	e.Static("/static", "public")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	})

	e.GET("/events", func(c echo.Context) error {
		events := make([]models.Event, 0)
		err := db.Where("public = ?", "1").Find(&events).Error

		if err != nil {
			panic(err)
		}

		return c.Render(http.StatusOK, "listevents.html", events)
	})

	e.GET("/places/:eventid", func(c echo.Context) error {
		eventid := c.Param("eventid")
		eventId, err := strconv.Atoi(eventid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum")
		}

		places := make([]models.Place, 0)
		err = db.Where("eventId = ?", eventId).Find(&places).Error

		if err != nil{
			panic(err)
		}

		event := new(models.Event)
		event.ID = uint(eventId)
		err = db.First(event).Error
		if err != nil || !event.Public {
			panic(err)
		}

		return c.Render(http.StatusOK, "listplacesforevent.html", struct {
			Event *models.Event
			Places []models.Place
		}{event, places})
	})

	e.GET("/participants/:eventid/:placeid", func(c echo.Context) error {
		placeid := c.Param("placeid")
		placeId, err := strconv.Atoi(placeid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum1")
		}

		eventid := c.Param("eventid")
		eventId, err := strconv.Atoi(eventid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum2")
		}

		place := new(models.Place)
		place.ID = uint(placeId)

		err = db.First(place).Error
		if err != nil{
			panic(err)
		}

		event := new(models.Event)
		event.ID = uint(eventId)

		err = db.First(event).Error
		if err != nil || !event.Public {
			panic(err)
		}

		participations := make([]models.Participation, 0)
		err = db.Where("place_id = ?", place.ID).Find(&participations).Error
		if err != nil {
			panic(err)
		}

		return c.Render(http.StatusOK, "participants.html", struct {
			Event *models.Event
			Place *models.Place
			Participants []models.Participation
		}{event, place, participations})
	})

	e.GET("/join/:eventid/:placeid", func(c echo.Context) error {
		placeid := c.Param("placeid")
		placeId, err := strconv.Atoi(placeid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum1")
		}

		eventid := c.Param("eventid")
		eventId, err := strconv.Atoi(eventid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum2")
		}

		place := new(models.Place)
		place.ID = uint(placeId)

		err = db.First(place).Error
		if err != nil{
			panic(err)
		}

		event := new(models.Event)
		event.ID = uint(eventId)

		err = db.First(event).Error
		if err != nil || !event.Public {
			panic(err)
		}

		return c.Render(http.StatusOK, "join.html", struct {
			Event *models.Event
			Place *models.Place
			EventId string
			PlaceId string
			Errors []string
		}{event, place,c.Param("eventid"), c.Param("placeid"), []string{}})
	})

	e.POST("/join/:eventid/:placeid", func(c echo.Context) error {
		Name := c.FormValue("Name")
		OM := c.FormValue("OM")


		placeid := c.Param("placeid")
		placeId, err := strconv.Atoi(placeid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum1")
		}

		eventid := c.Param("eventid")
		eventId, err := strconv.Atoi(eventid)
		if err != nil {
			return c.String(http.StatusBadRequest, "rossz formátum2")
		}

		place := new(models.Place)
		place.ID = uint(placeId)

		err = db.First(place).Error
		if err != nil{
			panic(err)
		}

		event := new(models.Event)
		event.ID = uint(eventId)

		err = db.First(event).Error
		if err != nil || !event.Public {
			panic(err)
		}

		student := new(models.Student)
		db.Where("name = ?", Name).Where("OM = ?", OM).First(student)

		errors := []string{}

		if student.ID == 0 {
			errors = append(errors, "Hibás név-om azonosító páros!")
		}

		if place.GetPeopleCount()>=place.PeopleCountLimit {
			errors = append(errors, "Ez a hely már megtelt!")
		}



		if len(errors) == 0 {
			participation := new(models.Participation)
			participation.StudentId = int(student.ID)
			participation.PlaceId = int(place.ID)

			ps := make([]models.Participation, 0)
			db.Where("student_id = ?", student.ID).First(&ps)
			existsInThisEvent := false

			for _, val := range ps {
				if val.Place().EventId == eventId {
					existsInThisEvent = true
					break
				}
			}

			if existsInThisEvent {
				errors = append(errors, "Ugyanabban az eseményben nem jelentkezhetsz különböző helyszínekre.")
			}else {
				err = db.Save(participation).Error
				if err != nil {
					panic(err)
				}
				return c.Redirect(http.StatusFound, "/participants/"+eventid+"/"+placeid)
			}
		}

		return c.Render(http.StatusOK, "join.html", struct {
			Event   *models.Event
			Place   *models.Place
			EventId string
			PlaceId string
			Errors  []string
		}{event, place, c.Param("eventid"), c.Param("placeid"), errors})

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

	adminClasses := admin.Group("/classes")
	adminClasses.GET("", getClasses)
	adminClasses.GET("/:id", getClass)
	adminClasses.PUT("/:id", putClass)
	adminClasses.POST("", postClass)
	adminClasses.DELETE("/:id", deleteClass)

	adminStudents := admin.Group("/students")
	adminStudents.GET("", getStudents)
	adminStudents.GET("/:id", getStudent)
	adminStudents.PUT("/:id", putStudent)
	adminStudents.POST("", postStudent)
	adminStudents.DELETE("/:id", deleteStudent)

	adminEvents := admin.Group("/events")
	adminEvents.GET("", getEvents)
	adminEvents.GET("/:id", getEvent)
	adminEvents.PUT("/:id", putEvent)
	adminEvents.POST("", postEvent)
	adminEvents.DELETE("/:id", deleteEvent)

	adminPlaces := admin.Group("/places")
	adminPlaces.GET("", getPlaces)
	adminPlaces.GET("/:id", getPlace)
	adminPlaces.PUT("/:id", putPlace)
	adminPlaces.POST("", postPlace)
	adminPlaces.DELETE("/:id", deletePlace)

	adminParticipations := admin.Group("/participations")
	adminParticipations.GET("", getParticipations)
	adminParticipations.GET("/:id", getParticipation)
	adminParticipations.PUT("/:id", putParticipation)
	adminParticipations.POST("", postParticipation)
	adminParticipations.DELETE("/:id", deleteParticipation)

	panic(e.Start(":8080"))
}
