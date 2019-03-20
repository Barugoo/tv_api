package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

//Valid constructs validates input fields and constructs error message
func (t *TV) Valid() error {
	errMsg := ""
	if !(t.ID > 0) {
		errMsg += fmt.Sprintf("id must be greater than %d. ", 0)
	}
	if !(len(t.Manufacturer) >= 3) {
		errMsg += fmt.Sprintf("manufacturer field must contain greater than or equal to %d symbols. ", 0)
	}
	if !(len(t.Model) >= 2) {
		errMsg += fmt.Sprintf("model field must contain greater than or equal to %d symbols. ", 2)
	}
	if !(t.Year >= 2010) {
		errMsg += fmt.Sprintf("year must be greater than or equal to %d.", 2010)
	}

	if len(errMsg) > 0 {
		return fmt.Errorf(errMsg)
	} else {
		return nil
	}
}

//NewTVService configures and initiates DB connection
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
	var input TV
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "create",
			"status": http.StatusBadRequest,
		}).Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "unable to unmarshhal json. field types mismatch?",
		})
		return
	}

	err = input.Valid()
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "create",
			"status": http.StatusBadRequest,
		}).Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": err.Error(),
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
		me, ok := err.(*mysql.MySQLError)
		if !ok || me.Number != 1062 {

			log.WithFields(log.Fields{
				"user":   r.RemoteAddr,
				"method": "create",
				"status": http.StatusInternalServerError,
			}).Error(me.Message)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseMsg{
				"status":  "error",
				"message": "database error",
			})
			return
		} else if me.Number == 1062 {

			log.WithFields(log.Fields{
				"user":   r.RemoteAddr,
				"method": "create",
				"status": http.StatusBadRequest,
			}).Error(me.Message)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseMsg{
				"status":  "error",
				"message": "dublicate key entry",
			})
			return
		}
	}

	log.WithFields(log.Fields{
		"user":   r.RemoteAddr,
		"method": "create",
		"status": http.StatusOK,
	}).Info()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseMsg{
		"status": "success",
		"data":   fmt.Sprintf("created record with id:%d", input.ID),
	})
}

func (s *TVService) Read(w http.ResponseWriter, r *http.Request) {
	result := &TV{}
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sqlStatement := "SELECT id, brand, manufacturer, model, year FROM tv WHERE id = ?"

	row := s.DB.QueryRow(sqlStatement, vars["id"])
	err = row.Scan(&result.ID, &result.Brand, &result.Manufacturer, &result.Model, &result.Year)
	if err != nil {

		if err.Error() == "sql: no rows in result set" {

			log.WithFields(log.Fields{
				"user":   r.RemoteAddr,
				"method": "read",
				"status": http.StatusNotFound,
			}).Error(err.Error())

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMsg{
				"status":  "error",
				"message": fmt.Sprintf("record with id:%s is not found", vars["id"]),
			})
			return
		}

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "read",
			"status": http.StatusInternalServerError,
		}).Error(err.Error())

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "database error",
		})
		return
	}

	log.WithFields(log.Fields{
		"user":   r.RemoteAddr,
		"method": "read",
		"status": http.StatusFound,
	}).Info()

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(result)
}

func (s *TVService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input TV
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "update",
			"status": http.StatusBadRequest,
		}).Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "unable to unmarshhal json. field types mismatch?",
		})
		return
	}

	err = input.Valid()
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "update",
			"status": http.StatusBadRequest,
		}).Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	sqlStatement := "UPDATE tv SET id = ?, brand = ?, manufacturer = ?, model = ?, year = ? WHERE id = ?"

	result, err := s.DB.Exec(
		sqlStatement,
		input.ID,
		input.Brand,
		input.Manufacturer,
		input.Model,
		input.Year,
		vars["id"],
	)
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "update",
			"status": http.StatusInternalServerError,
		}).Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "database error",
		})
		return
	}

	count, err := result.RowsAffected()
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "delete",
			"status": http.StatusInternalServerError,
		}).Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "database error",
		})
		return
	}

	if count == 0 {
		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "delete",
			"status": http.StatusNotFound,
		}).Error(fmt.Sprintf("nothing to update: record with id:%s is not found", vars["id"]))

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": fmt.Sprintf("nothing to update: record with id:%s is not found", vars["id"]),
		})
		return
	}

	log.WithFields(log.Fields{
		"user":   r.RemoteAddr,
		"method": "update",
		"status": http.StatusOK,
	}).Info()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseMsg{
		"status": "success",
		"data":   fmt.Sprintf("updated record with id:%s", vars["id"]),
	})

}

func (s *TVService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sqlStatement := "DELETE FROM tv WHERE id = ?"

	result, err := s.DB.Exec(
		sqlStatement,
		vars["id"],
	)
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "delete",
			"status": http.StatusInternalServerError,
		}).Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "database error",
		})
		return
	}

	count, err := result.RowsAffected()
	if err != nil {

		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "delete",
			"status": http.StatusInternalServerError,
		}).Error(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": "database error",
		})
		return
	}

	if count == 0 {
		log.WithFields(log.Fields{
			"user":   r.RemoteAddr,
			"method": "delete",
			"status": http.StatusNotFound,
		}).Error(fmt.Sprintf("nothing to delete: record with id:%s is not found", vars["id"]))

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMsg{
			"status":  "error",
			"message": fmt.Sprintf("nothing to delete: record wtih id:%s is not found", vars["id"]),
		})
		return
	}

	log.WithFields(log.Fields{
		"user":   r.RemoteAddr,
		"method": "delete",
		"status": http.StatusOK,
	}).Info()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseMsg{
		"status": "success",
		"data":   fmt.Sprintf("deleted record with id:%s", vars["id"]),
	})
}
