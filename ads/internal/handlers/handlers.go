package handlers

import (
	"github.com/gorilla/mux"
)

// ConfigureRouter configures this package http handlers for a given Gorilla Mux router
func ConfigureRouter(router *mux.Router) {
	router.HandleFunc(CreateURL, Create).Methods(CreateMethod)
	router.HandleFunc(UpdateURL, Update).Methods(UpdateMethod)
}
