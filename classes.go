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

	class.Id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	class.Scan(class.Id)

	return c.JSON(http.StatusOK, class)
}

func getClasses(c echo.Context) error {
	data, err := parsePaginationData(c)
	if err != nil {
		fmt.Println(err)
		return err
	}

	lst, err := models.ClassAPIGet(db, data._page, data._perPage, data._sortDir, data._sortField)
	if err != nil {
		fmt.Println(err, "models")
		return err
	}

	return c.JSON(http.StatusOK, lst)
}

func putClass(c echo.Context) error {
	var err error

	class := new(models.Class)
	if err := c.Bind(class); err != nil {
		fmt.Println("wut1", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}


	class.Id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("wut2", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = class.Update(db)
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

	err = class.Insert(db)
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
	err = class.Scan(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = class.Delete(db)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, class)
}