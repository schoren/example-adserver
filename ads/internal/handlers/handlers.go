package handlers

import (
	"github.com/gorilla/mux"
)

var AdServerBaseURL string

// ConfigureRouter configures this package http handlers for a given Gorilla Mux router
func ConfigureRouter(router *mux.Router) {
	router.HandleFunc(CreateURL, Create).Methods(CreateMethod)
	router.HandleFunc(UpdateURL, Update).Methods(UpdateMethod)
}
