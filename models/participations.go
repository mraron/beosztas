package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Participation struct {
	gorm.Model
	StudentId int  `gorm:"column:studentId"`
	PlaceId int  `gorm:"column:placeId"`
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

