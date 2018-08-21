package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"./models"
)

func getParticipation(c echo.Context) error {
	var err error

	object := new(models.Participation)

	object.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	db.Where("id = ?", object.ID).First(object)

	return c.JSON(http.StatusOK, object)
}

func getParticipations(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lst, err := models.ParticipationAPIGet(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cnt, _ := models.ParticipationAPIGetCount(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	c.Response().Header().Add("X-Total-Count", strconv.Itoa(cnt))

	return c.JSON(http.StatusOK, lst)
}

func putParticipation(c echo.Context) error {
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

func postParticipation(c echo.Context) error {
	var err error

	object := new(models.Participation)
	if err := c.Bind(object); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Create(object).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, object)
}

func deleteParticipation(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	object := new(models.Participation)
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

