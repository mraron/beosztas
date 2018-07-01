package models

import (
	"time"
	"database/sql"
)

type Event struct {
	Id int
	Name string
	Comment sql.NullString
	StartDate time.Time
	EndDate time.Time
	Public bool
}
