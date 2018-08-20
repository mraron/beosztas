package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Name string
	Comment string
	StartDate time.Time
	EndDate time.Time
	Public bool
}

func EventAPIGet(db *gorm.DB, _filters map[string]string, _page int, _perPage int, _sortDir string, _sortField string) ([]Event, error) {
	ans := make([]Event, 0)

	tmp := db.Order(_sortField+" "+_sortDir).Limit(_perPage).Offset(_perPage*(_page-1))
	for column, value := range _filters {
		tmp = tmp.Where(column+" like ?", "%"+value+"%")
	}

	err := tmp.Find(&ans).Error
	if err != nil {
		return nil, err
	}

	return ans, nil
}


