package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/parjinderpannu/garagesale/internal/mid"
	"github.com/parjinderpannu/garagesale/internal/platform/auth"
	"github.com/parjinderpannu/garagesale/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {

	app := web.NewApp(logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics())

	c := Check{DB: db}
	app.Handle(http.MethodGet, "/v1/health", c.Health)

	u := Users{DB: db, authenticator: authenticator}
	app.Handle(http.MethodGet, "/v1/users/token", u.Token)

	p := Product{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/v1/products", p.List, mid.Authenticate(authenticator))
	app.Handle(http.MethodPost, "/v1/products", p.Create, mid.Authenticate(authenticator))
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve, mid.Authenticate(authenticator))
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, mid.Authenticate(authenticator))
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, mid.Authenticate(authenticator),
		mid.HasRole(auth.RoleAdmin))

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator),
		mid.HasRole(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))

	return app
}
