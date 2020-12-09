package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/actions"
	"github.com/schoren/example-adserver/pkg/httputil"
)

// Define HTTP Method and URL
const (
	ListActiveMethod = http.MethodGet
	ListActiveURL    = "/active"
)

func NewListActive(a actions.ActiveLister) httputil.Handler {
	return &listActive{a}
}

type listActive struct {
	action actions.ActiveLister
}

func (h *listActive) Register(router *mux.Router) {
	router.HandleFunc(ListActiveURL, h.Handle).Methods(ListActiveMethod)
}

func (h *listActive) Handle(w http.ResponseWriter, r *http.Request) {
	ads, err := h.action.ListActive()
	if err != nil {
		log.Printf("Error getting active ads: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodedAds, _ := json.Marshal(ads)
	w.Header().Add("Content-Type", "application/json")
	w.Write(encodedAds)
}
