package actions_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/schoren/example-adserver/ads/internal/ads/actions"
	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockActiveAdGetter) ExpectGetActiveError(err error) {
	m.On("GetActive").Return([]types.Ad{}, err)
}

type listActiveFixtures struct {
	action   actions.ActiveLister
	adGetter *MockActiveAdGetter
	inst     *instrumentation.Mock
}

func (f listActiveFixtures) assertMockExpectations(t *testing.T) {
	f.adGetter.AssertExpectations(t)
	f.inst.AssertExpectations(t)
}

func listActiveSetup() listActiveFixtures {
	ag := new(MockActiveAdGetter)
	i := new(instrumentation.Mock)
	return listActiveFixtures{
		action:   actions.NewActiveLister(ag, i),
		adGetter: ag,
		inst:     i,
	}
}

func TestListActive(t *testing.T) {
	t.Parallel()

	var (
		exampleAds = []types.Ad{
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

		exampleError = fmt.Errorf("Some error with the data store!")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		f := listActiveSetup()
		f.adGetter.ExpectGetActiveSuccess(exampleAds)
		f.inst.ExpectOnStart()
		f.inst.ExpectOnComplete()

		ads, err := f.action.ListActive()

		assert.NoError(t, err)
		assert.Equal(t, exampleAds, ads)
		f.assertMockExpectations(t)
	})

	t.Run("AdsGetter error", func(t *testing.T) {
		t.Parallel()

		f := listActiveSetup()
		f.adGetter.ExpectGetActiveError(exampleError)
		f.inst.ExpectOnStart()
		f.inst.ExpectOnError(exampleError)
		f.inst.ExpectOnComplete()

		ads, err := f.action.ListActive()

		assert.Error(t, err)
		assert.True(t, errors.Is(err, exampleError))
		assert.Empty(t, ads)
		f.assertMockExpectations(t)
	})
}
