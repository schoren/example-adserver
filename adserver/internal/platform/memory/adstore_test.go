package memory_test

import (
	"errors"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/adserver/internal/platform/memory"
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

func TestAdStoreSet(t *testing.T) {
	as := memory.NewAdStore()

	as.Set(inmemoryExampleAdV1)

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV1, ad)
}

func TestAdStoreUpdate(t *testing.T) {
	as := memory.NewAdStore()

	as.Set(inmemoryExampleAdV1)

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV1, ad)

	as.Set(inmemoryExampleAdV2)

	ad, err = as.Get(inmemoryExampleAdV1.ID)
	assert.NoError(t, err)
	assert.Equal(t, inmemoryExampleAdV2, ad)
}

func TestAdStoreGetNotFound(t *testing.T) {
	as := memory.NewAdStore()

	ad, err := as.Get(inmemoryExampleAdV1.ID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, adstore.ErrNotFound))
	assert.Empty(t, ad)
}
