package models

type Place struct {
	Id int
	Name string
	Event *Event `db:"eventId"`
	PeopleCountLimit int
	ClassCountLimit int
	LowerClassLimit int
	UpperClassLimit int
}
