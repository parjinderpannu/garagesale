package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the signature used by all applicaiton handlers in this service
type Handler func(http.ResponseWriter, *http.Request) error

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
		if err := h(w, r); err != nil {
			// Log the error.
			a.log.Printf("ERROR : Unhandled error %v", err)
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
