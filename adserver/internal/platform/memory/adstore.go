package memory

import (
	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/pkg/types"
)

// AdStore is an in memory implementation of the AdStore
type AdStore struct {
	store map[int]types.Ad
}

// NewAdStore returns a ready to use InMemory adstore
func NewAdStore() *AdStore {
	return &AdStore{
		store: make(map[int]types.Ad),
	}
}

// Set an ad
func (s *AdStore) Set(ad types.Ad) {
	s.store[ad.ID] = ad
}

// Get an ad
func (s *AdStore) Get(id int) (types.Ad, error) {
	ad, ok := s.store[id]
	if !ok {
		return types.Ad{}, adstore.ErrNotFound
	}

	return ad, nil
}
