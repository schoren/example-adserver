package adstore_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSetter struct {
	mock.Mock
}

func (m *MockSetter) Set(ad types.Ad) {
	m.Called(ad)
}

type MockAdLister struct {
	mock.Mock
}

func (m *MockAdLister) List() ([]types.Ad, error) {
	args := m.Called()
	return args.Get(0).([]types.Ad), args.Error(1)
}

func warmupSetup() (*MockSetter, *MockAdLister) {
	store := new(MockSetter)
	adLister := new(MockAdLister)

	return store, adLister
}

func TestWarmup(t *testing.T) {
	t.Parallel()

	var (
		exampleAdList = []types.Ad{
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

		exampleListerError = fmt.Errorf("Some error")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		store, adLister := warmupSetup()

		adLister.On("List").Return(exampleAdList, nil)
		store.On("Set", exampleAdList[0]).Once()
		store.On("Set", exampleAdList[1]).Once()

		err := adstore.Warmup(store, adLister)

		assert.NoError(t, err)
		store.AssertExpectations(t)
		adLister.AssertExpectations(t)
	})

	t.Run("Lister error", func(t *testing.T) {
		t.Parallel()

		store, adLister := warmupSetup()

		adLister.On("List").Return([]types.Ad{}, exampleListerError)

		err := adstore.Warmup(store, adLister)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, exampleListerError))
		store.AssertExpectations(t)
		adLister.AssertExpectations(t)
	})
}
