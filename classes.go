package main

import (
	"github.com/labstack/echo"
	"fmt"
	"net/http"
	"./models"
	"strconv"
)

func getClass(c echo.Context) error {
	var err error

	class := new(models.Class)

	class.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	db.Where("id = ?", class.ID).First(class)

	return c.JSON(http.StatusOK, class)
}

func getClasses(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lst, err := models.ClassAPIGet(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cnt, _ := models.ClassAPIGetCount(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	c.Response().Header().Add("X-Total-Count", strconv.Itoa(cnt))

	return c.JSON(http.StatusOK, lst)
}

func putClass(c echo.Context) error {
	var err error

	class := new(models.Class)
	if err := c.Bind(class); err != nil {
		fmt.Println("wut1", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}


	class.ID, err = parseUint(c.Param("id"))
	if err != nil {
		fmt.Println("wut2", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Save(class).Error
	if err != nil {
		fmt.Println("wut3", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, class)
}

func postClass(c echo.Context) error {
	var err error

	class := new(models.Class)
	if err := c.Bind(class); err != nil {
		fmt.Println("wut1", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Create(class).Error
	if err != nil {
		fmt.Println("wut3", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, class)
}

func deleteClass(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	class := new(models.Class)
	err = db.Where("id = ?", id).First(class).Error

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Delete(class).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, class)
}