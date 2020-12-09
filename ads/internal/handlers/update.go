package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/schoren/example-adserver/ads/internal/actions"
	"github.com/schoren/example-adserver/pkg/httputil"
	"github.com/schoren/example-adserver/pkg/types"
)

// Define HTTP Method and URL
const (
	UpdateMethod = http.MethodPut
	UpdateURL    = "/{id}"
)

func NewUpdate(a actions.Updater) httputil.Handler {
	return &update{a}
}

type update struct {
	action actions.Updater
}

func (h *update) Register(router *mux.Router) {
	router.HandleFunc(UpdateURL, h.Handle).Methods(UpdateMethod)
}

func (h *update) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		log.Printf("Error getting ID from URL (%s): %v", r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid ID"}`))
		return
	}
	defer r.Body.Close()

	var req struct {
		ImageURL        string `json:"image_url"`
		ClickThroughURL string `json:"clickthrough_url"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid json"}`))
		return
	}

	ad := types.Ad{
		ID:              id,
		ImageURL:        req.ImageURL,
		ClickThroughURL: req.ClickThroughURL,
	}

	err = h.action.Update(ad)
	if err != nil {
		log.Printf("Error executing Update command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
