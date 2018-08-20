package models

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func SetDB(db_ *gorm.DB) {
	db=db_
}