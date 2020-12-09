package actions_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example-adserver/ads/internal/ads/actions"
	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/types"
)

type MockCreatePersister struct {
	mock.Mock
}

func (m *MockCreatePersister) Create(ad types.Ad) (types.Ad, error) {
	args := m.Called(ad)
	return args.Get(0).(types.Ad), args.Error(1)
}

func (m *MockCreatePersister) ExpectCreateSuccess(inputAd types.Ad, outputAd types.Ad) {
	m.On("Create", inputAd).Return(outputAd, nil)
}

func (m *MockCreatePersister) ExpectCreateError(inputAd types.Ad, err error) {
	m.On("Create", inputAd).Return(types.Ad{}, err)
}

type createFixtures struct {
	action    actions.Creator
	persister *MockCreatePersister
	notifier  *MockNotifier
	inst      *instrumentation.Mock
}

func (f createFixtures) assertMockExpectations(t *testing.T) {
	f.persister.AssertExpectations(t)
	f.notifier.AssertExpectations(t)
	f.inst.AssertExpectations(t)
}

func createSetup() createFixtures {
	p := new(MockCreatePersister)
	n := new(MockNotifier)
	i := new(instrumentation.Mock)
	c := actions.NewCreator(p, n, i)

	return createFixtures{c, p, n, i}
}

func TestCreat(t *testing.T) {
	t.Parallel()

	var (
		exampleAd = types.Ad{
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		examplePersistedAd = types.Ad{
			ID:              1,
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		exampleError = fmt.Errorf("Some error with the data store!")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		f := createSetup()

		f.persister.ExpectCreateSuccess(exampleAd, examplePersistedAd)
		f.notifier.ExpectAdUpdate(examplePersistedAd)
		f.inst.ExpectOnStart()
		f.inst.ExpectOnComplete()

		ad, err := f.action.Create(exampleAd)

		assert.NoError(t, err)
		assert.Equal(t, examplePersistedAd, ad)
		f.assertMockExpectations(t)
	})

	t.Run("persister error", func(t *testing.T) {
		f := createSetup()
		f.persister.ExpectCreateError(exampleAd, exampleError)

		f.inst.ExpectOnStart()
		f.inst.ExpectOnError(exampleError)
		f.inst.ExpectOnComplete()

		ad, err := f.action.Create(exampleAd)

		assert.True(t, errors.Is(err, exampleError))
		assert.Equal(t, types.Ad{}, ad)
		f.assertMockExpectations(t)
	})
}
