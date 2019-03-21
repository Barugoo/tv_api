package main

import (
	"database/sql"

	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CaseResponse
type CR map[string]interface{}

type Case struct {
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
	Path   string
	Query  string
	Status int
	Result interface{}
	Body   interface{}
}

var (
	client = &http.Client{Timeout: time.Second}
)

func PrepareTestApis(db *sql.DB) {
	qs := []string{
		`DROP TABLE IF EXISTS tv;`,

		`CREATE TABLE tv (
			id int(11) NOT NULL,
			brand varchar(255),
			manufacturer varchar(255) NOT NULL,
			model varchar(255) NOT NULL,
			year int(11) NOT NULL,
			PRIMARY KEY (id)
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	}

	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}

func CleanupTestApis(db *sql.DB) {
	qs := []string{
		`DROP TABLE IF EXISTS tv;`,
	}
	for _, q := range qs {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
	}
}
