package adstore

import (
	"fmt"

	"github.com/schoren/example-adserver/types"
)

var ErrNotFound = fmt.Errorf("Ad not found")

// InMemory is an in memory implementation of the AdStore
type InMemory struct {
	store map[int]types.Ad
}

// NewInMemory returns a ready to use InMemory adstore
func NewInMemory() *InMemory {
	return &InMemory{
		store: make(map[int]types.Ad),
	}
}

// Set an ad
func (s *InMemory) Set(ad types.Ad) {
	s.store[ad.ID] = ad
}

// Get an ad
func (s *InMemory) Get(id int) (types.Ad, error) {
	ad, ok := s.store[id]
	if !ok {
		return types.Ad{}, ErrNotFound
	}

	return ad, nil
}
