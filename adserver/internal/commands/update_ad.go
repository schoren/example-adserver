package commands

import (
	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/types"
)

type UpdateAdPayload struct {
	Ad types.Ad
}

// UpdateAdCommand updates the adstore with new or updated ads
type UpdateAdCommand struct {
	AdStore adstore.Setter
}

func (c *UpdateAdCommand) Execute(payload UpdateAdPayload) error {
	c.AdStore.Set(payload.Ad)
	return nil
}
