package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/commands"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type MockAdStore struct {
	mock.Mock
}

func (m *MockAdStore) Get(id int) (types.Ad, error) {
	args := m.Called(id)
	return args.Get(0).(types.Ad), args.Error(1)
}

func (m *MockAdStore) ExpectGetSuccess(id int, ad types.Ad) {
	m.On("Get", id).Return(ad, nil)
}

func (m *MockAdStore) ExpectGetError(id int) {
	m.On("Get", id).Return(types.Ad{}, serveExampleAdStoreError)
}

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
