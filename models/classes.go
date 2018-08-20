package models

import (
	"github.com/jinzhu/gorm"
)

type Class struct {
	gorm.Model
	Index int
	Name string
}


func ClassAPIGet(db *gorm.DB, _filters map[string]string, _page int, _perPage int, _sortDir string, _sortField string) ([]Class, error) {
	ans := make([]Class, 0)

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
