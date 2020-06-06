package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/types"
)

// Define HTTP Method and URL
const (
	CreateMethod = http.MethodPost
	CreateURL    = "/"
)

// Creater handles the creation of a given ad
type Creater interface {
	Execute(commands.CreatePayload) (types.Ad, error)
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

	persistedAd, err := CreateCommand.Execute(payload)
	if err != nil {
		log.Printf("Error executing Create command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s/%d", AdServerURL, persistedAd.ID))
	w.WriteHeader(http.StatusCreated)
}
