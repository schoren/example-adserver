package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/types"
)

// Define HTTP Method and URL
const (
	CreateMethod = http.MethodPost
	CreateURL    = "/create"
)

// Creater handles the creation of a given ad
type Creater interface {
	Execute(commands.CreatePayload) error
}

var CreateCommand Creater

type createRequest struct {
	ImageURL        string `json:"image_url"`
	ClickThroughURL string `json:"clickthrough_url"`
}

// Create handles HTTP requests for creating ads
func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ad createRequest
	err := json.NewDecoder(r.Body).Decode(&ad)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid json"}`))
		return
	}

	payload := commands.CreatePayload{
		Ad: types.Ad{
			ImageURL:        ad.ImageURL,
			ClickThroughURL: ad.ClickThroughURL,
		},
	}

	err = CreateCommand.Execute(payload)
	if err != nil {
		log.Printf("Error executing Create command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
