package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(logger *log.Logger, db *sqlx.DB) http.Handler {

	p := Product{DB: db, Log: logger}

	return http.HandlerFunc(p.List)
}
