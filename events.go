package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"./models"
)

func getEvent(c echo.Context) error {
	var err error

	object := new(models.Event)

	object.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	db.Where("id = ?", object.ID).First(object)

	return c.JSON(http.StatusOK, object)
}

func getEvents(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lst, err := models.EventAPIGet(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cnt, _ := models.EventAPIGetCount(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	c.Response().Header().Add("X-Total-Count", strconv.Itoa(cnt))

	return c.JSON(http.StatusOK, lst)
}

func putEvent(c echo.Context) error {
	var err error

	object := new(models.Event)
	if err := c.Bind(object); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	object.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Save(object).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, object)
}

func postEvent(c echo.Context) error {
	var err error

	object := new(models.Event)
	if err := c.Bind(object); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Create(object).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, object)
}

func deleteEvent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	object := new(models.Event)
	err = db.Where("id = ?", id).First(object).Error

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Delete(object).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, object)
}
