package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/parjinderpannu/garagesale/internal/platform/auth"
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
func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(ctx, p.DB)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

// Retrieve gives a single Product.
func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(ctx, p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "getting product %q", id)
		}
	}

	return web.Respond(ctx, w, prod, http.StatusOK)
}

// Create decode a JSON  document from a POST request
// and create a new Product
func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("auth claims not in context")
	}

	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(ctx, p.DB, claims, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)
}

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL.
func (p *Product) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	if err := product.Update(ctx, p.DB, claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete gives a single Product.
func (p *Product) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

	if err := product.Delete(ctx, p.DB, id); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body. The full model is returned to the caller.
func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrap(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(ctx, p.DB, ns, productID, time.Now())
	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product.
func (p *Product) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(ctx, p.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}
