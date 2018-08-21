package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"./models"
)

func getPlace(c echo.Context) error {
	var err error

	object := new(models.Place)

	object.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	db.Where("id = ?", object.ID).First(object)

	return c.JSON(http.StatusOK, object)
}

func getPlaces(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lst, err := models.PlaceAPIGet(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cnt, _ := models.PlaceAPIGetCount(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	c.Response().Header().Add("X-Total-Count", strconv.Itoa(cnt))

	return c.JSON(http.StatusOK, lst)
}

func putPlace(c echo.Context) error {
	var err error

	object := new(models.Place)
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

func postPlace(c echo.Context) error {
	var err error

	object := new(models.Place)
	if err := c.Bind(object); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Create(object).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, object)
}

func deletePlace(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	object := new(models.Place)
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
