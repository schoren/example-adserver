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

var (
	warmupExampleAdList = []types.Ad{
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

	warmupExampleListerError = fmt.Errorf("Some error")
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

func TestWarmupOK(t *testing.T) {
	store, adLister := warmupSetup()

	adLister.On("List").Return(warmupExampleAdList, nil)
	store.On("Set", warmupExampleAdList[0]).Once()
	store.On("Set", warmupExampleAdList[1]).Once()

	err := adstore.Warmup(store, adLister)

	assert.NoError(t, err)
	store.AssertExpectations(t)
	adLister.AssertExpectations(t)
}

func TestWarmupListerError(t *testing.T) {
	store, adLister := warmupSetup()

	adLister.On("List").Return([]types.Ad{}, warmupExampleListerError)

	err := adstore.Warmup(store, adLister)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, warmupExampleListerError))
	store.AssertExpectations(t)
	adLister.AssertExpectations(t)
}
