package models

import "github.com/jmoiron/sqlx"

var db *sqlx.DB

func SetDB(db_ *sqlx.DB) {
	db=db_
}