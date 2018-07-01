package models

type Participation struct {
	Id int
	Student *Student `db:"studentId"`
	Place *Place `db:"PlaceId"`
}