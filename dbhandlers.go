package main

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrDBError      = errors.New("db error")
	ErrDublicateKey = errors.New("dublicate key")
	ErrUpdate       = errors.New("same data")
)

func (s *TVService) DBCreate(input TV) error {

	sqlStatement := "INSERT INTO tv (id, brand, manufacturer, model, year) VALUES (?, ?, ?, ?, ?)"

	_, err := s.DB.Exec(
		sqlStatement,
		input.ID,
		input.Brand,
		input.Manufacturer,
		input.Model,
		input.Year,
	)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok || me.Number != 1062 {
			return ErrDBError
		} else if me.Number == 1062 {
			return ErrDublicateKey
		}
	}
	return nil
}

func (s *TVService) DBRead(id int) (*TV, error) {
	result := &TV{}

	sqlStatement := "SELECT id, brand, manufacturer, model, year FROM tv WHERE id = ?"

	row := s.DB.QueryRow(sqlStatement, id)
	err := row.Scan(&result.ID, &result.Brand, &result.Manufacturer, &result.Model, &result.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		} else {
			return nil, ErrDBError
		}
	}
	return result, nil
}

func (s *TVService) DBUpdate(input TV, id int) error {

	sqlUpdateStatement := "UPDATE tv SET id = ?, brand = ?, manufacturer = ?, model = ?, year = ? WHERE id = ?"
	sqlCheckStatement := "SELECT id FROM tv WHERE id = ?"

	dummy := &TV{}
	row := s.DB.QueryRow(sqlCheckStatement, id)
	err := row.Scan(&dummy.ID)

	if err != nil {

		if err == sql.ErrNoRows {
			return ErrNotFound
		} else {
			return ErrDBError
		}
	}

	result, err := s.DB.Exec(
		sqlUpdateStatement,
		input.ID,
		input.Brand,
		input.Manufacturer,
		input.Model,
		input.Year,
		id,
	)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok || me.Number != 1062 {
			return ErrDBError
		} else if me.Number == 1062 {
			return ErrDublicateKey
		}
	}

	count, err := result.RowsAffected()
	if err != nil {
		return ErrDBError
	}

	if count == 0 {
		return ErrUpdate
	}

	return nil
}

func (s *TVService) DBDelete(id int) error {

	sqlStatement := "DELETE FROM tv WHERE id = ?"

	result, err := s.DB.Exec(
		sqlStatement,
		id,
	)
	if err != nil {
		return ErrDBError
	}
	count, err := result.RowsAffected()
	if count == 0 {
		return ErrNotFound
	} else if err != nil {
		return ErrDBError
	}
	return nil
}
