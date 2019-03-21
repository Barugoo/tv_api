package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

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

	fmt.Println("starting server at :" + service.Cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+service.Cfg.ServerPort, r))
}
