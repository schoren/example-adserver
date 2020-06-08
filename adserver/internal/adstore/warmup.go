package adstore

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/types"
)

type AdLister interface {
	List() ([]types.Ad, error)
}

// Warmup a given store with the results of the given AdLister
func Warmup(store Setter, lister AdLister) error {
	ads, err := lister.List()
	if err != nil {
		return fmt.Errorf("Error listing ads for warmup: %w", err)
	}

	for _, ad := range ads {
		store.Set(ad)
	}

	return nil
}
