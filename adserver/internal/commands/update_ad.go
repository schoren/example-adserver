package commands

import (
	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/types"
)

type UpdateAdPayload struct {
	Ad types.Ad
}

// UpdateAd updates the adstore with new or updated ads
type UpdateAd struct {
	AdStore adstore.Setter
}

func NewUpdateAd(adStore adstore.Setter) *UpdateAd {
	return &UpdateAd{AdStore: adStore}
}

func (c *UpdateAd) Execute(payload UpdateAdPayload) error {
	c.AdStore.Set(payload.Ad)
	return nil
}
