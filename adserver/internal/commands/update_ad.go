package commands

import "github.com/schoren/example-adserver/types"

type AdStoreSetter interface {
	Set(types.Ad)
}

type UpdateAdPayload struct {
	Ad types.Ad
}

// UpdateAdCommand updates the adstore with new or updated ads
type UpdateAdCommand struct {
	AdStore AdStoreSetter
}

func (c *UpdateAdCommand) Execute(payload UpdateAdPayload) error {
	c.AdStore.Set(payload.Ad)
	return nil
}
