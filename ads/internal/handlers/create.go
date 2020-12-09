package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/actions"
	"github.com/schoren/example-adserver/pkg/httputil"
	"github.com/schoren/example-adserver/pkg/types"
)

// Define HTTP Method and URL
const (
	CreateMethod = http.MethodPost
	CreateURL    = "/"
)

func NewCreate(a actions.Creator, baseURL string) httputil.Handler {
	return &create{
		action:  a,
		baseURL: baseURL,
	}
}

type create struct {
	action  actions.Creator
	baseURL string
}

func (h *create) Register(router *mux.Router) {
	router.HandleFunc(CreateURL, h.Handle).Methods(CreateMethod)
}

func (h *create) Handle(w http.ResponseWriter, r *http.Request) {
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

	ad, err = h.action.Create(ad)
	if err != nil {
		log.Printf("Error executing Create command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s/%d", h.baseURL, ad.ID))
	w.WriteHeader(http.StatusCreated)
}
