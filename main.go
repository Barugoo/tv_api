package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	DSN        string
	Driver     string
	ServerPort string
}

func main() {
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

	fmt.Println("starting server at :" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}