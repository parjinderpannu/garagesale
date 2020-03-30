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
}

// NewApp Knows how to construct internal state for an App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

// Handle connects a method and URL pattern to a
// particular application handler.
func (a *App) Handle(method, pattern string, h Handler) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {

			// Log the error.
			a.log.Printf("ERROR : %+v", err)

			// Respond to the error.
			if err := RespondError(w, err); err != nil {
				a.log.Printf("ERROR : %v", err)
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
