package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/adserver/internal/actions"
	"github.com/schoren/example-adserver/adserver/internal/renderer"
	"github.com/schoren/example-adserver/pkg/httputil"
)

const (
	ServeMethod = http.MethodGet
	ServeURL    = "/{id}"
)

type Server interface {
	Serve(int) (renderer.Renderer, error)
}

func NewServer(a actions.Server) httputil.Handler {
	return &server{a}
}

type server struct {
	action actions.Server
}

func (h *server) Register(router *mux.Router) {
	router.HandleFunc(ServeURL, h.Handle).Methods(ServeMethod)
}

func (h *server) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		log.Printf("Error getting ID from URL (%s): %v", r.RequestURI, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	renderer, err := h.action.Serve(id)
	if err != nil {
		log.Printf("Error Serving ad with ID %d: %v", id, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Fprint(w, renderer.Render())
}
