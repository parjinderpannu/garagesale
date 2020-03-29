package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/parjinderpannu/garagesale/internal/platform/web"
	"github.com/parjinderpannu/garagesale/internal/product"
)

// Product has handler method for dealing with Products.
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List is a HTTP Handler for returning a list of Products.
func (p *Product) List(w http.ResponseWriter, r *http.Request) {

	p.Log.Println("cmd/sales-api/internal/handler/product.go testing")

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

// Retrieve gives a single Product.
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusOK); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}

// Create decode a JSON  document from a POST request
// and create a new Product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) {

	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying db", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("error responding", err)
		return
	}
}
