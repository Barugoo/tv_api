package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

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
	return &TVService{DB: db, Cfg: cfg}, nil
}

//ResponseService makes a response
func ResponseService(w http.ResponseWriter, user string, method string, status int, msg string, responseObject interface{}) {
	log.WithFields(log.Fields{
		"user":   user,
		"method": method,
		"status": status,
	}).Error(msg)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseObject)
	return
}

//Valid validates input fields and constructs error message
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

func (s *TVService) Create(w http.ResponseWriter, r *http.Request) {
	var input TV
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		msg := "unable to unmarshhal json. field types mismatch?"
		ResponseService(w, r.RemoteAddr, "create", http.StatusBadRequest, msg, ResponseMsg{
			"status":  "error",
			"message": msg,
		})
		return
	}

	err = input.Valid()
	if err != nil {
		msg := err.Error()
		ResponseService(w, r.RemoteAddr, "create", http.StatusBadRequest, msg, ResponseMsg{
			"status":  "error",
			"message": msg,
		})
		return
	}

	err = s.DBCreate(input)

	if err != nil {
		switch err {
		case ErrDBError:
			msg := "internal server error"
			ResponseService(w, r.RemoteAddr, "create", http.StatusInternalServerError, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrDublicateKey:
			msg := fmt.Sprintf("duplicate key entry, id:%d", input.ID)
			ResponseService(w, r.RemoteAddr, "create", http.StatusBadRequest, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return
		}
	}

	msg := fmt.Sprintf("created record with id:%d", input.ID)
	ResponseService(w, r.RemoteAddr, "create", http.StatusOK, msg, ResponseMsg{
		"status":  "success",
		"message": msg,
	})
	return
}

func (s *TVService) Read(w http.ResponseWriter, r *http.Request) {
	result := &TV{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	result, err = s.DBRead(id)

	if err != nil {
		switch err {
		case ErrDBError:
			msg := "internal server error"
			ResponseService(w, r.RemoteAddr, "read", http.StatusInternalServerError, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrNotFound:
			msg := fmt.Sprintf("record with id:%d is not found", id)
			ResponseService(w, r.RemoteAddr, "read", http.StatusNotFound, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return
		}
	}
	ResponseService(w, r.RemoteAddr, "read", http.StatusOK, "OK!", result)
	return
}

func (s *TVService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input TV
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&input)
	if err != nil {
		msg := "unable to unmarshhal json. field types mismatch?"
		ResponseService(w, r.RemoteAddr, "update", http.StatusBadRequest, msg, ResponseMsg{
			"status":  "error",
			"message": msg,
		})
		return
	}

	err = input.Valid()
	if err != nil {
		msg := err.Error()
		ResponseService(w, r.RemoteAddr, "update", http.StatusBadRequest, msg, ResponseMsg{
			"status":  "error",
			"message": msg,
		})
		return
	}

	err = s.DBUpdate(input, id)

	if err != nil {
		switch err {
		case ErrDBError:
			msg := "internal server error"
			ResponseService(w, r.RemoteAddr, "update", http.StatusInternalServerError, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrNotFound:
			msg := fmt.Sprintf("nothing to update: record with id:%d is not found", id)
			ResponseService(w, r.RemoteAddr, "update", http.StatusNotFound, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrUpdate:
			msg := fmt.Sprintf("same data provided: record id:%d wasn't updated", id)
			ResponseService(w, r.RemoteAddr, "update", http.StatusBadRequest, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrDublicateKey:
			msg := fmt.Sprintf("duplicate key entry, id:%d", input.ID)
			ResponseService(w, r.RemoteAddr, "update", http.StatusBadRequest, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return
		}

	}

	msg := fmt.Sprintf("updated record with id:%d", id)
	ResponseService(w, r.RemoteAddr, "update", http.StatusOK, msg, ResponseMsg{
		"status":  "success",
		"message": msg,
	})
	return
}

func (s *TVService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err = s.DBDelete(id)

	if err != nil {
		switch err {
		case ErrDBError:
			msg := "internal server error"
			ResponseService(w, r.RemoteAddr, "delete", http.StatusInternalServerError, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return

		case ErrNotFound:
			msg := fmt.Sprintf("nothing to delete: record with id:%d is not found", id)
			ResponseService(w, r.RemoteAddr, "delete", http.StatusNotFound, msg, ResponseMsg{
				"status":  "error",
				"message": msg,
			})
			return
		}
	}

	msg := fmt.Sprintf("deleted record with id:%d", id)
	ResponseService(w, r.RemoteAddr, "delete", http.StatusOK, msg, ResponseMsg{
		"status":  "success",
		"message": msg,
	})
	return
}
