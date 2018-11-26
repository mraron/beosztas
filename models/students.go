package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sort"
	"strconv"
)

type Student struct {
	gorm.Model
	Name string
	OM string `gorm:"column:OM"`
	ClassId int `json:"class_id"`
	Count int `gorm:"-"`
}

func StudentAPIGet(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) ([]Student, error) {
	ans := make([]Student, 0)

	tmp := db.Model(&Student{})

	if _sortField != "Count" {
		tmp = tmp.Order(_sortField + " " + _sortDir).Limit(_perPage).Offset(_perPage*(_page-1))
	}

	filter_column := -1

	for column, value := range _filters {
		fmt.Println(column,"=>", value)
		if column == "class_id" {
			tmp = tmp.Where(column+" = ?", value)
		}else if column == "Count" {
			kaki, err := strconv.Atoi(fmt.Sprintf("%v", value))
			fmt.Println(err)
			filter_column = int(kaki)
		} else {
			tmp = tmp.Where(column+" like ?", fmt.Sprintf("%%%v%%", value))
		}
	}

	err := tmp.Find(&ans).Error
	if err != nil {
		return nil, err
	}

	fmt.Println(filter_column,"!!!!!!!!")
	ans2 := make([]Student, 0, len(ans))
	for ind := range ans {
		ans[ind].Init()
		if filter_column != -1 {
			if ans[ind].Count == filter_column {
				ans2 = append(ans2, ans[ind])
			}
		}
	}

	if filter_column != -1 {
		ans2, ans = ans, ans2
	}

	if _sortField == "Count" {
		sort.Slice(ans, func(i, j int) bool {
			if _sortDir == "ASC" {
				return ans[i].Count < ans[j].Count
			}

			return ans[i].Count > ans[j].Count
		})

		lim := _perPage*_page
		if lim > len(ans) {
			lim = len(ans)
		}

		return ans[_perPage*(_page-1):lim], nil
	}


	return ans, nil
}

func StudentAPIGetCount(db *gorm.DB, _filters map[string]interface{}, _page int, _perPage int, _sortDir string, _sortField string) (int, error) {
	lst, err := StudentAPIGet(db, _filters, 1, 1000000, _sortDir, _sortField)
	return len(lst), err
}

func (s *Student) Class() *Class {
	class := new(Class)
	class.ID = uint(s.ClassId)
	db.First(class)

	return class
}

func (s *Student) Init() {
	mp := make(map[int]int)
	mp[10] = 2
	mp[9] = 1
	mp[8] = 1
	mp[7] = 1

	ans := 0

	ossz := make([]Participation, 0)

	db.Where("student_id=?", s.ID).Find(&ossz)
	fmt.Println(s.ID, len(ossz), "!!!!!!!!!!!")
	for _, val := range ossz {

		ans+=mp[int(val.Place().EventId)]
	}
	fmt.Println(ans)
	s.Count = ans
}


