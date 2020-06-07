package adstore_test

import (
	"errors"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
)

var (
	inmemoryExampleAdV1 = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org",
	}

	inmemoryExampleAdV2 = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif?updated",
		ClickThroughURL: "http://example.org/?updated",
	}
)

func TestInMemorySet(t *testing.T) {
	as := adstore.NewInMemory()

	as.Set(inmemoryExampleAdV1)

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV1, ad)
}

func TestInMemoryUpdate(t *testing.T) {
	as := adstore.NewInMemory()

	as.Set(inmemoryExampleAdV1)

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV1, ad)

	as.Set(inmemoryExampleAdV2)

	ad, err = as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV2, ad)
}

func TestInMemoryGetNotFound(t *testing.T) {
	as := adstore.NewInMemory()

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, adstore.ErrNotFound))
	assert.Empty(t, ad)
}
