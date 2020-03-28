package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
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

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling result", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing result", err)
	}

}
