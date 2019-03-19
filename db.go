package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type TV struct {
	ID           int64  `json:"id"`
	Brand        string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Year         string `json:"year"`
}

func (t *TV) Valid() bool {
	year, err := strconv.Atoi(t.Year)
	if err != nil {
		return false
	}
	return t.ID > 0 && len(t.Manufacturer) >= 3 && len(t.Model) >= 2 && year >= 2010
}

type TVService struct {
	DB *sql.DB
}

func NewTVService(cfg Config) (*TVService, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &TVService{DB: db}, nil
}

func (s *TVService) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	var input TV
	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	if !input.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	sqlStatement := "INSERT INTO tv (id, brand, manufacturer, model, year) VALUES (?, ?, ?, ?, ?)"

	_, err = s.DB.Exec(
		sqlStatement,
		input.ID,
		input.Brand,
		input.Manufacturer,
		input.Model,
		input.Year,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "unable to communicate with database",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   nil,
	})
}

func (s *TVService) Read(w http.ResponseWriter, r *http.Request) {
	result := &TV{}
	vars := mux.Vars(r)

	sqlStatement := "SELECT id, brand, manufacturer, model, year FROM tv WHERE id = ?"

	row := s.DB.QueryRow(sqlStatement, vars["id"])
	err := row.Scan(&result.ID, &result.Brand, &result.Manufacturer, &result.Model, &result.Year)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "not found",
		})
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (s *TVService) Update(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	var input TV
	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	if !input.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "bad request",
		})
		return
	}

	vars := mux.Vars(r)

	sqlStatement := "UPDATE tv SET id = ?, brand = ?, manufacturer = ?, model = ?, year = ? WHERE id = ?"

	_, err = s.DB.Exec(
		sqlStatement,
		input.ID,
		input.Brand,
		input.Manufacturer,
		input.Model,
		input.Year,
		vars["id"],
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "unable to communicate with database",
		})
		return
	}
	http.Redirect(w, r, "/"+strconv.Itoa(int(input.ID)), http.StatusFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   nil,
	})

}
func (s *TVService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sqlStatement := "DELETE FROM tv WHERE id = ?"

	_, err := s.DB.Exec(
		sqlStatement,
		vars["id"],
	)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"message": "unable to communicate with database",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   nil,
	})

}
