package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

// Handler is the signature used by all applicaiton handlers in this service
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// App is the entrypoint for all web applications.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

// NewApp constructs an App to handle a set of routes.
// Any Middleware provided will be ran for every request.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}
}

// Handle connects a method and URL pattern to a
// particular application handler.
func (a *App) Handle(method, pattern string, h Handler) {

	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {

		// Create a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := h(ctx, w, r); err != nil {
			// Log the error.
			a.log.Printf("ERROR : Unhandled error %v", err)
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
