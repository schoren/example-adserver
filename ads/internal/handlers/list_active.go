package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/schoren/example-adserver/types"
)

// Define HTTP Method and URL
const (
	ListActiveMethod = http.MethodGet
	ListActiveURL    = "/"
)

// ActiveLister handles the creation of a given ad
type ActiveLister interface {
	Execute() ([]types.Ad, error)
}

var ListActiveCommand ActiveLister

// ListActive handles HTTP requests for listing currently active ads
func ListActive(w http.ResponseWriter, r *http.Request) {
	ads, err := ListActiveCommand.Execute()
	if err != nil {
		log.Printf("Error getting active ads: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodedAds, _ := json.Marshal(ads)
	w.Header().Add("Content-Type", "application/json")
	w.Write(encodedAds)
}
