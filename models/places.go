package models

import (
	"database/sql/driver"
	"strconv"
	"errors"
)

type Place struct {
	Id int
	Name string
	Event *Event `db:"eventId"`
	PeopleCountLimit int
	ClassCountLimit int
	LowerClassLimit int
	UpperClassLimit int
}

func (s Place) Value() (driver.Value, error) {
	return driver.Value(s.Id), nil
}

func (s *Place) Scan(value interface{}) error {
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

	row := db.QueryRow("SELECT * FROM places WHERE id=?", id)
	if err = row.Scan(&s.Id, &s.Name, &s.Event, &s.PeopleCountLimit, &s.ClassCountLimit, &s.LowerClassLimit, &s.UpperClassLimit); err != nil {
		return err
	}

	return nil
}
