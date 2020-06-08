package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	listActiveExampleAds = []types.Ad{
		{
			ID:              1,
			ImageURL:        "https://example.org/1.png",
			ClickThroughURL: "https://example.org/1.html",
		},
		{
			ID:              2,
			ImageURL:        "https://example.org/2.png",
			ClickThroughURL: "https://example.org/2.html",
		},
	}

	listActiveExamplePersisterError = fmt.Errorf("Some error with the data store!")
)

type MockActiveAdGetter struct {
	mock.Mock
}

func (m *MockActiveAdGetter) GetActive() ([]types.Ad, error) {
	args := m.Called()
	return args.Get(0).([]types.Ad), args.Error(1)
}

func (m *MockActiveAdGetter) ExpectGetActiveSuccess(returnAds []types.Ad) {
	m.On("GetActive").Return(returnAds, nil)
}

func (m *MockActiveAdGetter) ExpectGetActiveError() {
	m.On("GetActive").Return([]types.Ad{}, listActiveExamplePersisterError)
}

func listActiveSetup() (*commands.ListActive, *MockActiveAdGetter) {
	adsGetter := new(MockActiveAdGetter)
	cmd := commands.NewListActive(adsGetter)

	return cmd, adsGetter
}

func TestListActiveOK(t *testing.T) {
	cmd, adsGetter := listActiveSetup()
	adsGetter.ExpectGetActiveSuccess(listActiveExampleAds)

	ads, err := cmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, listActiveExampleAds, ads)
}

func TestListActiveAdsGetterError(t *testing.T) {
	cmd, adsGetter := listActiveSetup()
	adsGetter.ExpectGetActiveError()

	ads, err := cmd.Execute()

	assert.Error(t, err)
	assert.True(t, errors.Is(err, listActiveExamplePersisterError))
	assert.Empty(t, ads)
}
