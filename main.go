package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/parjinderpannu/garagesale/schema"
)

func main() {

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	// it happens when func main is done
	defer log.Println("main : Completed")

	// =========================================================================
	// Setup Dependencies

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations", err)
		}
		log.Println("Migration Successful")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed data", err)
		}
		log.Println("Seed data successful")
		return
	}

	// =========================================================================
	// Start API Service

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

//Product is an item we sell.
type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ListProducts is a HTTP Handler for returning a list of Products.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Product{
		{Name: "Comic Books", Cost: 50, Quantity: 42},
		{Name: "McDonalds Toys", Cost: 75, Quantity: 120},
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshalling result", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing result", err)
	}

}
