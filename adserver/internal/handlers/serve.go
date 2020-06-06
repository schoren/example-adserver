package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/adserver/internal/renderer"
)

const (
	ServeMethod = http.MethodGet
	ServeURL    = "/{id}"
)

type Server interface {
	Execute(int) (renderer.Renderer, error)
}

var ServeCommand Server

// Serve handles HTTP requests for ad serving
func Serve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		log.Printf("Error getting ID from URL (%s): %v", r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	renderer, err := ServeCommand.Execute(id)
	if err != nil {
		log.Printf("Error Serving ad with ID %d: %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprint(w, renderer.Render())
	w.WriteHeader(http.StatusOK)
}
