package mid

import (
	"log"
	"net/http"
	"time"

	"github.com/parjinderpannu/garagesale/internal/platform/web"
)

// Logger will log a line for every request.
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			start := time.Now()
			// Run the handler chain and catch any propagated error.
			err := before(w, r)

			log.Printf(
				"%s %s (%v)",
				r.Method, r.URL.Path, time.Since(start),
			)

			// Return the error to be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
