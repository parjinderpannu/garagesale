package mid

import (
	"log"
	"net/http"

	"github.com/parjinderpannu/garagesale/internal/platform/web"
)

// Logger will log a line for every request.
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			// Run the handler chain and catch any propagated error.
			err := before(w, r)

			log.Printf(
				"%s %s",
				r.Method, r.URL.Path,
			)

			// Return the error to be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
