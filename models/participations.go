package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Participation struct {
	gorm.Model
	StudentId int `json:"student_id"`
	PlaceId int `json:"place_id"`
	ClassId uint `json:"class_id"`
}

func ParticipationAPIGet(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) ([]Participation, error) {
	ans := make([]Participation, 0)

	tmp := db.Order(_sortField+" "+_sortDir).Limit(_perPage).Offset(_perPage*(_page-1))

	for column, value := range _filters {
		tmp = tmp.Where(column+" like ?", fmt.Sprintf("%%%v%%", value))
	}

	err := tmp.Find(&ans).Error
	if err != nil {
		return nil, err
	}

	for ind, _ := range ans {
		ans[ind].ClassId = ans[ind].Class().ID
	}

	return ans, nil
}

func ParticipationAPIGetCount(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) (int, error) {
	ans := 0

	tmp := db.Model(&Participation{}).Order(_sortField+" "+_sortDir)

	for column, value := range _filters {
		tmp = tmp.Where(column+" like ?", fmt.Sprintf("%%%v%%", value))
	}

	err := tmp.Count(&ans).Error
	if err != nil {
		return -1, err
	}

	return ans, nil
}

func (p *Participation) Student() *Student {
	student := new(Student)
	student.ID = uint(p.StudentId)

	db.First(student)
	return student
}


func (p *Participation) Place() *Place {
	place := new(Place)
	place.ID = uint(p.PlaceId)

	db.First(place)
	return place
}

func (p *Participation) Class() *Class {
	class := new(Class)
	class.ID = uint(p.Student().ClassId)

	db.First(class)
	return class
}
