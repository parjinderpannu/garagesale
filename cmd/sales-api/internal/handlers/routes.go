package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/parjinderpannu/garagesale/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger)

	p := Product{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)

	return app
}
