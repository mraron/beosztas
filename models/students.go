package models

import (
	"github.com/jinzhu/gorm"
)

type Student struct {
	gorm.Model
	Name string
	OM string `gorm:"column:OM"`
	ClassId int `gorm:"column:classId"`
}

func StudentAPIGet(db *gorm.DB, _filters map[string]string, _page int, _perPage int, _sortDir string, _sortField string) ([]Student, error) {
	ans := make([]Student, 0)

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

func (s *Student) Class() *Class {
	class := new(Class)
	class.ID = uint(s.ClassId)
	db.First(class)

	return class
}


