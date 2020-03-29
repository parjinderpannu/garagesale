package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/parjinderpannu/garagesale/internal/platform/web"
	"github.com/parjinderpannu/garagesale/internal/product"
	"github.com/pkg/errors"
)

// Product has handler method for dealing with Products.
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List is a HTTP Handler for returning a list of Products.
func (p *Product) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(p.DB)
	if err != nil {
		return errors.Wrap(err, "getting product list")
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve gives a single Product.
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		return errors.Wrapf(err, "getting product %q", id)
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decode a JSON  document from a POST request
// and create a new Product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return errors.Wrap(err, "decoding new product")
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "creating new product")
	}

	return web.Respond(w, prod, http.StatusCreated)
}
