package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"./models"
)

func getStudent(c echo.Context) error {
	var err error

	student := new(models.Student)

	student.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	db.Where("id = ?", student.ID).First(student)

	return c.JSON(http.StatusOK, student)
}

func getStudents(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	lst, err := models.StudentAPIGet(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	cnt, _ := models.StudentAPIGetCount(db, data._filters, data._page, data._perPage, data._sortDir, data._sortField)
	c.Response().Header().Add("X-Total-Count", strconv.Itoa(cnt))

	return c.JSON(http.StatusOK, lst)
}

func putStudent(c echo.Context) error {
	var err error

	student := new(models.Student)
	if err := c.Bind(student); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	student.ID, err = parseUint(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Save(student).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, student)
}

func postStudent(c echo.Context) error {
	var err error

	student := new(models.Student)
	if err := c.Bind(student); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Create(student).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, student)
}

func deleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	student := new(models.Student)
	err = db.Where("id = ?", id).First(student).Error

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = db.Delete(student).Error
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, student)
}
