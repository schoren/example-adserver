package commands

import (
	"fmt"

	"github.com/schoren/example-adserver/ads/internal/types"
)

type CreatePayload struct {
	Ad types.Ad
}

// CreatePersister allows to persist an ad to a data store
type CreatePersister interface {
	Create(types.Ad) (types.Ad, error)
}

// Create is a command used to create a new Ad
type Create struct {
	Persister CreatePersister
	Notifier  Notifier
}

// Execute the Create command with the given payload
func (c Create) Execute(data CreatePayload) error {
	ad, err := c.Persister.Create(data.Ad)
	if err != nil {
		return fmt.Errorf("Persister.Create error when creating ad: %w", err)
	}

	c.Notifier.AdUpdate(ad)

	return nil
}
