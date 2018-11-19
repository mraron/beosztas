package main

import (
	"./models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"github.com/SaidinWoT/timespan"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	HOSTNAME = ""
	DATABASENAME = ""
	PORT = ""
	ADMIN_USER = ""
	ADMIN_PASSWORD = ""
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

	db, err = gorm.Open("sqlite3", DATABASENAME)
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
	_filters map[string]interface{}
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
	res._filters = make(map[string]interface{})

	err = json.Unmarshal([]byte(_filters), &res._filters)
	return res, nil
}


func parseUint(str string) (uint, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(res), nil
}

func currentStudent(c echo.Context) (*models.Student, error) {
	var (
		u   *models.Student = &models.Student{}
		err error
	)

	storage, err := session.Get("student", c)
	if err != nil {
		panic(err)
	}

	if _, ok := storage.Values["id"]; !ok {
		return nil, nil
	}

	db.Where("ID=?", storage.Values["id"]).First(u)

	return u, err
}

func internalError(c echo.Context, err error, msg string) error {
	c.Logger().Print("internal error:", err)
	return c.Render(http.StatusInternalServerError, "error.html", msg)
}

func loadConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	config := struct {
		Databasename string
		Hostname string
		Port string
		Admin_user string
		Admin_password string
	}{"","", "", "", ""}

	d := json.NewDecoder(f)
	err = d.Decode(&config)
	if err != nil {
		panic(err)
	}

	DATABASENAME = config.Databasename
	HOSTNAME = config.Hostname
	PORT = config.Port
	ADMIN_USER = config.Admin_user
	ADMIN_PASSWORD = config.Admin_password
}

func main() {
	log.Print("loading config")
	loadConfig()

	log.Print("connecting to db")
	connectToDB()
	models.SetDB(db)

	st := new(models.Student)
	db.First(st)

	log.Print("populating echo web framework")
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

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("titkosdolog"))))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user, err := currentStudent(c)
			fmt.Println(user, err)
			if err != nil {
				return internalError(c, err, "belső hiba")
			}

			c.Set("student", user)
			fmt.Println(c.Get("student"))
			return next(c)
		}
	})

	e.Static("/static", "public")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	})

	e.GET("/login", func(c echo.Context) error {
		if c.Get("student") == nil || c.Get("student").(*models.Student)==nil  {
			return c.Render(http.StatusOK, "login.html", nil)
		}

		return c.Redirect(http.StatusFound, "/")
	})

	e.POST("/login", func(c echo.Context) error {
		if c.Get("student") != nil && c.Get("student").(*models.Student)!=nil  {
			return c.Redirect(http.StatusFound, "/")
		}

		name, om := c.FormValue("name"), c.FormValue("om")
		ember := models.Student{}
		db.Where("name = ?", name).Where("OM = ?", om).First(&ember)

		if ember.ID==0 {
			return c.Render(http.StatusOK, "login.html", "nincs ilyen név-om pár.")
		}

		storage, _ := session.Get("student", c)
		storage.Values["id"] = ember.ID

		if err := storage.Save(c.Request(), c.Response()); err != nil {
			return internalError(c, err, "Belső hiba #2")
		}

		c.Set("student", ember)


		return c.Redirect(http.StatusFound, "/my/show/")
	})

	e.GET("/logout", func(c echo.Context) error {
		if c.Get("student") == nil || c.Get("student").(*models.Student)==nil  {
			return c.Redirect(http.StatusFound, "/login")
		}

		storage, _ := session.Get("student", c)
		storage.Options.MaxAge = -1
		storage.Values["id"] = -1

		if err := storage.Save(c.Request(), c.Response()); err != nil {
			return internalError(c, err, "Belső hiba")
		}

		return c.Redirect(http.StatusFound, "/")
	})

	e.GET("/my/show/", func(c echo.Context) error {
		fmt.Println( c.Get("student") == nil, "%!3453453454354353453453453")
		if c.Get("student") == nil || c.Get("student").(*models.Student)==nil {
			return c.Redirect(http.StatusFound, "/login")
		}

		s := c.Get("student").(*models.Student)

		nev, om := s.Name, s.OM

		ember := models.Student{}
		db.Where("name = ?", nev).Where("OM = ?", om).Find(&ember)

		if ember.ID==0 {
			return c.Render(http.StatusOK, "my.html", "nincs ilyen név-om pár.")
		}

		fmt.Println(ember, "!!")

		ans := make([]models.Participation, 0)
		db.Where("student_id = ?", ember.ID).Find(&ans)
		return c.Render(http.StatusOK, "myshow.html", struct {
			Ans []models.Participation
			Nev string
			OM string
		}{ans, nev, om})
	})

	e.POST("/my/action", func(c echo.Context) error {  //@TODO lul it's shiet
		action := c.FormValue("action")
		id := c.FormValue("id")
		nev := c.FormValue("nev")
		om := c.FormValue("om")

		fmt.Println("action", id, action, nev, om)

		student := models.Student{}

		db.Where("name = ?", nev).Where("OM = ?", om).Find(&student)

		part := models.Participation{}

		db.Where("ID = ?", id).Find(&part)

		if uint(part.StudentId) == student.ID {
			if action=="del" {
				db.Where("ID = ?", id).Delete(models.Participation{})
			}else {
				part.Validated = true
				db.Save(&part)
			}
		}

		return c.Redirect(http.StatusFound, c.FormValue("redirect"))
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
		err = db.Where("event_id = ?", eventId).Find(&places).Error

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

	/*e.GET("/join/:eventid/:placeid", func(c echo.Context) error {
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
	})*/

	e.POST("/join/:eventid/:placeid", func(c echo.Context) error {
		if c.Get("student") == nil || c.Get("student").(*models.Student)==nil  {
			return c.Redirect(http.StatusFound, "/login")
		}
		s := c.Get("student").(*models.Student)

		Name, OM := s.Name, s.OM

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
			db.Where("student_id = ?", student.ID).Find(&ps)
			ok := true

			span := timespan.New(place.Event().StartDate, place.Event().EndDate.Sub(place.Event().StartDate))
			fmt.Println(span)
			fmt.Println(place.Event())
			fmt.Println("???????????????????????????????")
			for _, val := range ps {
				akt := timespan.New(val.Place().Event().StartDate, val.Place().Event().EndDate.Sub(val.Place().Event().StartDate))
				fmt.Println(akt)
				fmt.Println(val.Place().Event())
				huh, van := akt.Intersection(span)
				fmt.Println(huh, van)
				if  van {
					ok = false
					errors = append(errors, fmt.Sprintf("Ugyanabban az időpontban nem lehetsz két helyen is! Az ütköző helyszín %q és esemény %q.", val.Place().Name, val.Place().Event().Name))
					break
				}
			}

			if !ok {

			}else {
				err = db.Save(participation).Error
				if err != nil {
					panic(err)
				}
				return c.Redirect(http.StatusFound, fmt.Sprintf("/my/show/"))
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
		if username == ADMIN_USER && password == ADMIN_PASSWORD {
			return true, nil
		}
		return false, nil
	}))

	admin.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "admin.html", HOSTNAME)
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

	panic(e.Start(":"+PORT))
}
