package models

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Place struct {
	gorm.Model
	Name             string `gorm:"column:name"`
	Location string
	PeopleCountLimit int `gorm:"column:peopleCountLimit"`
	EventId int  `gorm:"column:eventId"`
}

func PlaceAPIGet(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) ([]Place, error) {
	ans := make([]Place, 0)

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

func PlaceAPIGetCount(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) (int, error) {
	ans := 0

	tmp := db.Model(&Place{}).Order(_sortField+" "+_sortDir)

	for column, value := range _filters {
		tmp = tmp.Where(column+" like ?", fmt.Sprintf("%%%v%%", value))
	}

	err := tmp.Count(&ans).Error
	if err != nil {
		return -1, err
	}

	return ans, nil
}

func (p *Place) GetPeopleCount() int {
	count := 0
	db.Model(new(Participation)).Where("place_id = ?", p.ID).Count(&count)
	return count
}