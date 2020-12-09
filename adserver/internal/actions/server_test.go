package actions_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/actions"
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

func setupServer() (actions.Server, *MockAdStore) {
	mockAdStore := new(MockAdStore)
	action := actions.NewServer(mockAdStore)

	return action, mockAdStore
}

func TestServeOK(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		action, mockAdStore := setupServer()
		mockAdStore.ExpectGetSuccess(serveExampleAd.ID, serveExampleAd)

		renderer, err := action.Serve(serveExampleAd.ID)

		assert.NoError(t, err)
		assert.Equal(t, renderer.Render(), serveExampleRenderedAd)
		mockAdStore.AssertExpectations(t)
	})

	t.Run("AdStore error", func(t *testing.T) {
		t.Parallel()

		action, mockAdStore := setupServer()
		mockAdStore.ExpectGetError(serveExampleAd.ID)

		renderer, err := action.Serve(serveExampleAd.ID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, serveExampleAdStoreError))
		assert.Empty(t, renderer)
		mockAdStore.AssertExpectations(t)
	})
}
