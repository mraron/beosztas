package models

import (
	"database/sql/driver"
	"strconv"
	"errors"
)

type Participation struct {
	Id int
	Student *Student `db:"studentId"`
	Place *Place `db:"p	laceId"`
}

func (s Participation) Value() (driver.Value, error) {
	return driver.Value(s.Id), nil
}

func (s *Participation) Scan(value interface{}) error {
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

	row := db.QueryRow("SELECT * FROM participations WHERE id=?", id)
	if err = row.Scan(&s.Id, &s.Student, &s.Place); err != nil {
		return err
	}

	return nil
}
