package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

//Config struct
type Config struct {
	DSN        string
	Driver     string
	ServerPort string
}

//TV struct
type TV struct {
	ID           int64  `json:"id"`
	Brand        string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
}

//ResponseMsg is used for encoding a response
type ResponseMsg map[string]interface{}

//TVService handles DB connection and
type TVService struct {
	DB     *sql.DB
	Router *mux.Router
	Cfg    Config
}
