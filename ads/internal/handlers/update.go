package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/types"
)

// Define HTTP Method and URL
const (
	UpdateMethod = http.MethodPut
	UpdateURL    = "/{id}"
)

// Updater handles the creation of a given ad
type Updater interface {
	Execute(commands.UpdatePayload) error
}

var UpdateCommand Updater

type updateRequest struct {
	ImageURL        string `json:"image_url"`
	ClickThroughURL string `json:"clickthrough_url"`
}

// Update handles HTTP requests for creating ads
func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		log.Printf("Error getting ID from URL (%s): %v", r.RequestURI, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid ID"}`))
		return
	}
	defer r.Body.Close()

	var ad updateRequest
	err = json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid json"}`))
		return
	}

	payload := commands.UpdatePayload{
		Ad: types.Ad{
			ID:              id,
			ImageURL:        ad.ImageURL,
			ClickThroughURL: ad.ClickThroughURL,
		},
	}

	err = UpdateCommand.Execute(payload)
	if err != nil {
		log.Printf("Error executing Update command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
