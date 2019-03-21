package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// CaseResponse
type CR map[string]interface{}

type Case struct {
	Method string
	Path   string
	Status int
	Body   interface{}
	Result interface{}
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

		`INSERT INTO tv (id, brand, manufacturer, model, year) VALUES
		  (1,	'Bravia',	'Sony',	'HX929',	2011);`,
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

func TestApis(t *testing.T) {
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/api")
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	PrepareTestApis(db)

	defer CleanupTestApis(db)

	cfg := Config{
		DSN:        "root:1234@tcp(localhost:3306)/api",
		Driver:     "mysql",
		ServerPort: "8082",
	}

	service, err := NewTVService(cfg)
	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/").Subrouter()
	tv := api.PathPrefix("/tv/").Subrouter()

	tv.HandleFunc("/new", service.Create).Methods(http.MethodPost)
	tv.HandleFunc("/{id}", service.Read).Methods(http.MethodGet)
	tv.HandleFunc("/{id}", service.Update).Methods(http.MethodPut)
	tv.HandleFunc("/{id}", service.Delete).Methods(http.MethodDelete)

	ts := httptest.NewServer(r)
	cases := []Case{
		Case{
			Method: "GET",
			Path:   "/666",
			Status: http.StatusNotFound,
			Result: CR{
				"status":  "error",
				"message": "record with id:666 is not found",
			},
		},
		Case{
			Method: "GET",
			Path:   "/1",
			Status: http.StatusOK,
			Result: CR{
				"id":           1,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "HX929",
				"year":         2011,
			},
		},
		Case{
			Method: "POST",
			Path:   "/new",
			Status: http.StatusOK,
			Body: CR{
				"id":           2,
				"brand":        "Smart",
				"manufacturer": "Philips",
				"model":        "8000",
				"year":         2012,
			},
			Result: CR{
				"status":  "success",
				"message": "created record with id:2",
			},
		},
		Case{
			Method: "POST",
			Path:   "/new",
			Status: http.StatusOK,
			Body: CR{
				"id":           3,
				"manufacturer": "LG-TV",
				"model":        "LA8600",
				"year":         2013,
			},
			Result: CR{
				"status":  "success",
				"message": "created record with id:3",
			},
		},
		Case{
			Method: "POST",
			Path:   "/new",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           4,
				"brand":        "Colour Television",
				"manufacturer": "GoldStar Co.",
				"model":        "CB-14A80",
				"year":         1992,
			},
			Result: CR{
				"status":  "error",
				"message": "year must be greater than or equal to 2010.",
			},
		},
		Case{
			Method: "POST",
			Path:   "/new",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           3,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         2014,
			},
			Result: CR{
				"status":  "error",
				"message": "duplicate key entry, id:3",
			},
		},
		Case{
			Method: "POST",
			Path:   "/new",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           4,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         "2014",
			},
			Result: CR{
				"status":  "error",
				"message": "unable to unmarshhal json. field types mismatch?",
			},
		},
		Case{
			Method: "PUT",
			Path:   "/2",
			Status: http.StatusOK,
			Body: CR{
				"id":           2,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         2014,
			},
			Result: CR{
				"status":  "success",
				"message": "updated record with id:2",
			},
		},
		Case{
			Method: "PUT",
			Path:   "/2",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           1,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         2014,
			},
			Result: CR{
				"status":  "error",
				"message": "duplicate key entry, id:1",
			},
		},
		Case{
			Method: "PUT",
			Path:   "/2",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           1,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         "2014",
			},
			Result: CR{
				"status":  "error",
				"message": "unable to unmarshhal json. field types mismatch?",
			},
		},
		Case{
			Method: "PUT",
			Path:   "/2",
			Status: http.StatusBadRequest,
			Body: CR{
				"id":           2,
				"brand":        "Colour Television",
				"manufacturer": "GoldStar Co.",
				"model":        "CB-14A80",
				"year":         1992,
			},
			Result: CR{
				"status":  "error",
				"message": "year must be greater than or equal to 2010.",
			},
		},
		Case{
			Method: "PUT",
			Path:   "/666",
			Status: http.StatusNotFound,
			Body: CR{
				"id":           666,
				"brand":        "Bravia",
				"manufacturer": "Sony",
				"model":        "X950B",
				"year":         2014,
			},
			Result: CR{
				"status":  "error",
				"message": "nothing to update: record with id:666 is not found",
			},
		},
		Case{
			Method: "DELETE",
			Path:   "/1",
			Status: http.StatusOK,
			Result: CR{
				"status":  "success",
				"message": "deleted record with id:1",
			},
		},
		Case{
			Method: "DELETE",
			Path:   "/1",
			Status: http.StatusNotFound,
			Result: CR{
				"status":  "error",
				"message": "nothing to delete: record with id:1 is not found",
			},
		},
	}
	runCases(t, ts, db, cases)
}

func runCases(t *testing.T, ts *httptest.Server, db *sql.DB, cases []Case) {
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s", idx, item.Method, item.Path)

		if item.Method == "" || item.Method == http.MethodGet {
			req, err = http.NewRequest(item.Method, ts.URL+"/api/tv"+item.Path, nil)
		} else {
			data, err := json.Marshal(item.Body)
			if err != nil {
				panic(err)
			}
			reqBody := bytes.NewReader(data)
			req, err = http.NewRequest(item.Method, ts.URL+"/api/tv"+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if item.Status == 0 {
			item.Status = http.StatusOK
		}

		if resp.StatusCode != item.Status {
			t.Fatalf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		data, err := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("[%s] results not match\nGot : %#v\nWant: %#v", caseName, result, expected)
			continue
		}
	}

}
