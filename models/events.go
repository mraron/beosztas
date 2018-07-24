package models

import (
	"time"
	"database/sql"
	"database/sql/driver"
	"strconv"
	"errors"
)

type Event struct {
	Id int
	Name string
	Comment sql.NullString
	StartDate time.Time
	EndDate time.Time
	Public bool
}

func (s Event) Value() (driver.Value, error) {
	return driver.Value(s.Id), nil
}

func (s *Event) Scan(value interface{}) error {
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

	row := db.QueryRow("SELECT * FROM events WHERE id=?", id)
	if err = row.Scan(&s.Id, &s.Name, &s.Comment, &s.StartDate, &s.EndDate, &s.Public); err != nil {
		return err
	}

	return nil
}
