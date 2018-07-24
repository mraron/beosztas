package models

import (
	"database/sql/driver"
	"strconv"
	"errors"
	"github.com/jmoiron/sqlx"
	"fmt"
)

type Class struct {
	Id int
	Index int
	Name string
}

func (s Class) Value() (driver.Value, error) {
	return driver.Value(s.Id), nil
}

func (s *Class) Scan(value interface{}) error {
	if value == nil {
		return errors.New("Can't scan student from nil")																																																										
	}

	var (
		id int
		err error
	)

	switch value.(type) {
	case int64:
		id = int(value.(int64))
	case int:
		id = value.(int)
	case []uint8:
		if id, err = strconv.Atoi(string(value.([]uint8))); err != nil {
			return err
		}
	}

	row := db.QueryRow("SELECT * FROM classes WHERE id=?", id)
	if err = row.Scan(&s.Id, &s.Index, &s.Name); err != nil {
		return err
	}

	return nil
}


func (pr Class) Delete(db *sqlx.DB) (err error) {
	_, err = db.Exec("DELETE FROM classes WHERE id=$1", pr.Id)
	return
}

func (pr Class) Update(db *sqlx.DB) (err error) {
	_, err = db.Exec("UPDATE classes SET \"index\"=$1, name=$2 WHERE id=$3", pr.Index, pr.Name, pr.Id)
	return
}

func (pr *Class) Insert(db *sqlx.DB) error {
	res, err := db.Exec("INSERT INTO classes (\"index\", name) VALUES ($1,$2)", pr.Index, pr.Name)

	id_, err := res.LastInsertId()

	if err != nil {
		return err
	}

	pr.Id = int(id_)

	return err
}


func getClasses(db *sqlx.DB, query string, args ...interface{}) (ans []Class, err error) {
	var rows *sqlx.Rows

	rows, err = db.Queryx(query, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		class := Class{}
		err = rows.Scan(&class.Id, &class.Index, &class.Name)
		if err != nil {
			break
		}

		ans = append(ans, class)
	}

	if err = rows.Err(); err != nil {
		return
	}

	rows.Close()

	return
}


func ClassAPIGet(db *sqlx.DB, _page int, _perPage int, _sortDir string, _sortField string) ([]Class, error) {
	lst, err := getClasses(db, fmt.Sprintf("SELECT * FROM classes ORDER BY %s %s LIMIT %d OFFSET %d ", _sortField, _sortDir, _perPage, _perPage*(_page-1))) //@TODO: SQL Injection!!!
	if err != nil {
		return nil, err
	}

	return lst, nil
}
