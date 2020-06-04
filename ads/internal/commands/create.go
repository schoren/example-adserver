package commands

import (
	"fmt"

	"github.com/schoren/example/ads/ads/internal/types"
)

type CreatePayload struct {
	Ad types.Ad
}

// Persister allows to persist an ad to a data store
type Persister interface {
	Create(types.Ad) (types.Ad, error)
}

// Notifier can propagate events to other components of the system
type Notifier interface {
	AdCreated(types.Ad)
}

// Create is a command used to create a new Ad
type Create struct {
	Persister Persister
	Notifier  Notifier
}

// Execute the Create command with the given payload
func (c Create) Execute(data CreatePayload) error {
	ad, err := c.Persister.Create(data.Ad)
	if err != nil {
		return fmt.Errorf("Persister.Create error when creating ad: %w", err)
	}

	c.Notifier.AdCreated(ad)

	return nil
}
