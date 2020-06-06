package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/commands"
	"github.com/schoren/example-adserver/types"
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

func setupServe() (commands.ServeCommand, *MockAdStore) {
	mockAdStore := new(MockAdStore)
	cmd := commands.ServeCommand{
		AdStore: mockAdStore,
	}

	return cmd, mockAdStore
}

func TestServeCommandOK(t *testing.T) {
	cmd, mockAdStore := setupServe()
	mockAdStore.ExpectGetSuccess(serveExampleAd.ID, serveExampleAd)

	renderer, err := cmd.Execute(serveExampleAd.ID)

	assert.NoError(t, err)
	assert.Equal(t, renderer.Render(), serveExampleRenderedAd)
	mockAdStore.AssertExpectations(t)
}

func TestServeCommandAdStoreError(t *testing.T) {
	cmd, mockAdStore := setupServe()
	mockAdStore.ExpectGetError(serveExampleAd.ID)

	renderer, err := cmd.Execute(serveExampleAd.ID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, serveExampleAdStoreError))
	assert.Empty(t, renderer)
	mockAdStore.AssertExpectations(t)
}
