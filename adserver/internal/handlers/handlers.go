package handlers

import (
	"github.com/gorilla/mux"
)

// ConfigureRouter configures this package http handlers for a given Gorilla Mux router
func ConfigureRouter(router *mux.Router) {
	router.HandleFunc(ServeURL, Serve).Methods(ServeMethod)
}
