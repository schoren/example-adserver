package rest

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/schoren/example-adserver/types"
	"gopkg.in/resty.v1"
)

type AdLister struct {
	client  *resty.Client
	baseURL string
}

func NewAdLister(adServicBaseURL string) *AdLister {
	return &AdLister{
		client:  resty.New(),
		baseURL: adServicBaseURL,
	}
}

func (l *AdLister) List() ([]types.Ad, error) {
	resp, err := l.client.R().Get(fmt.Sprintf("%s/ads/active", l.baseURL))
	if err != nil {
		return []types.Ad{}, fmt.Errorf("Failed to get ads: %v", err)
	}

	ads := []types.Ad{}
	log.Println("Resp was", string(resp.Body()))
	err = json.Unmarshal(resp.Body(), &ads)
	if err != nil {
		return []types.Ad{}, fmt.Errorf("Failed to decode ads response: %v", err)
	}

	return ads, nil
}
