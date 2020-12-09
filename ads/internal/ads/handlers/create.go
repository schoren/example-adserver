package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/pkg/types"
)

// Define HTTP Method and URL
const (
	CreateMethod = http.MethodPost
	CreateURL    = "/"
)

var CreateCommand commands.Creator

// Create handles HTTP requests for creating ads
func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req struct {
		ImageURL        string `json:"image_url"`
		ClickThroughURL string `json:"clickthrough_url"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid json"}`))
		return
	}

	ad := types.Ad{
		ImageURL:        req.ImageURL,
		ClickThroughURL: req.ClickThroughURL,
	}

	ad, err = CreateCommand.Execute(ad)
	if err != nil {
		log.Printf("Error executing Create command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s/%d", AdServerBaseURL, ad.ID))
	w.WriteHeader(http.StatusCreated)
}
