package commands

import (
	"fmt"

	"github.com/schoren/example-adserver/ads/internal/types"
)

type UpdatePayload struct {
	Ad types.Ad
}

// UpdatePersister allows to persist an ad to a data store
type UpdatePersister interface {
	Update(types.Ad) error
}

// Update is a command used to Update a new Ad
type Update struct {
	Persister UpdatePersister
	Notifier  Notifier
}

// Execute the Update command with the given payload
func (c Update) Execute(data UpdatePayload) error {
	err := c.Persister.Update(data.Ad)
	if err != nil {
		return fmt.Errorf("Persister.Update error when updating ad: %w", err)
	}

	c.Notifier.AdUpdate(data.Ad)

	return nil
}
