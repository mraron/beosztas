package models

import (
	"database/sql"
	"database/sql/driver"

	"errors"
	"strconv"
)

type Student struct {
	Id int
	Name string
	OM sql.NullString
	Class *Class `db:"classId"`
}

func (s Student) Value() (driver.Value, error) {
	return driver.Value(s.Id), nil
}

func (s *Student) Scan(value interface{}) error {
	if value == nil {
		return errors.New("Can't scan student from nil")
	}

	var (
		id int
		err error
	)

	switch value.(type) {
	case int64:
		id = int(value.(int64))
	case int:
		id = value.(int)
	case []uint8:
		if id, err = strconv.Atoi(string(value.([]uint8))); err != nil {
			return err
		}
	}

	row := db.QueryRow("SELECT * FROM students WHERE id=?", id)
	if err = row.Scan(&s.Id, &s.Name, &s.OM, &s.Class	); err != nil {
		return err
	}

	return nil
}
