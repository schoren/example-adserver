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

type MockUpdatePersister struct {
	mock.Mock
}

func (m *MockUpdatePersister) Update(ad types.Ad) error {
	args := m.Called(ad)
	return args.Error(0)
}

func (m *MockUpdatePersister) ExpectUpdateSuccess(inputAd types.Ad) {
	m.On("Update", inputAd).Return(nil)
}

func (m *MockUpdatePersister) ExpectUpdateError(err error) {
	m.On("Update", mock.Anything).Return(err)
}

type updateFixtures struct {
	action    actions.Updater
	persister *MockUpdatePersister
	notifier  *MockNotifier
	inst      *instrumentation.Mock
}

func (f updateFixtures) assertMockExpectations(t *testing.T) {
	f.persister.AssertExpectations(t)
	f.notifier.AssertExpectations(t)
	f.inst.AssertExpectations(t)
}

func setupUpdate() updateFixtures {
	p := new(MockUpdatePersister)
	n := new(MockNotifier)
	i := new(instrumentation.Mock)
	a := actions.NewUpdater(p, n, i)

	return updateFixtures{a, p, n, i}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	var (
		exampleAd = types.Ad{
			ID:              1,
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		exampleError = fmt.Errorf("Some error with the data store!")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		f := setupUpdate()

		f.persister.ExpectUpdateSuccess(exampleAd)
		f.notifier.ExpectAdUpdate(exampleAd)
		f.inst.ExpectOnStart()
		f.inst.ExpectOnComplete()

		err := f.action.Update(exampleAd)

		assert.NoError(t, err)
		f.assertMockExpectations(t)
	})

	t.Run("Persist error", func(t *testing.T) {
		t.Parallel()

		f := setupUpdate()
		f.persister.ExpectUpdateError(exampleError)
		f.inst.ExpectOnStart()
		f.inst.ExpectOnError(exampleError)
		f.inst.ExpectOnComplete()

		err := f.action.Update(exampleAd)

		assert.True(t, errors.Is(err, exampleError))
		f.assertMockExpectations(t)
	})
}
