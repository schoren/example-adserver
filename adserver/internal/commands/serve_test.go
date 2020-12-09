package actions_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/commands"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	serveExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org",
	}
	serveExampleRenderedAd   = `<a href="http://example.org"><img src="http://example.org/img.gif"></a>`
	serveExampleAdStoreError = fmt.Errorf("Some datastore error")
)

func setupServe() (*commands.Serve, *MockAdStore) {
	mockAdStore := new(MockAdStore)
	cmd := commands.NewServe(mockAdStore)

	return cmd, mockAdStore
}

func TestServeOK(t *testing.T) {
	cmd, mockAdStore := setupServe()
	mockAdStore.ExpectGetSuccess(serveExampleAd.ID, serveExampleAd)

	renderer, err := cmd.Execute(serveExampleAd.ID)

	assert.NoError(t, err)
	assert.Equal(t, renderer.Render(), serveExampleRenderedAd)
	mockAdStore.AssertExpectations(t)
}

func TestServeAdStoreError(t *testing.T) {
	cmd, mockAdStore := setupServe()
	mockAdStore.ExpectGetError(serveExampleAd.ID)

	renderer, err := cmd.Execute(serveExampleAd.ID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, serveExampleAdStoreError))
	assert.Empty(t, renderer)
	mockAdStore.AssertExpectations(t)
}
