package models

import "database/sql"

type Student struct {
	Id int
	Name string
	Password sql.NullString
	Email sql.NullString
	Class *Class `db:"classId"`
	ActivationCode string
}


