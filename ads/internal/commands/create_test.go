package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example/ads/ads/internal/commands"
	"github.com/schoren/example/ads/ads/internal/types"
)

type MockPersister struct {
	mock.Mock
}

func (m *MockPersister) Create(ad types.Ad) (types.Ad, error) {
	args := m.Called(ad)
	return args.Get(0).(types.Ad), args.Error(1)
}

func (m *MockPersister) ExpectCreateSuccess(inputAd types.Ad, outputAd types.Ad) {
	m.On("Create", inputAd).Return(outputAd, nil)
}

func (m *MockPersister) ExpectCreateError(inputAd types.Ad, err error) {
	m.On("Create", inputAd).Return(types.Ad{}, err)
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) AdCreated(inputAd types.Ad) {
	m.Called(inputAd)
}

func (m *MockNotifier) ExpectAdCreated(inputAd types.Ad) {
	m.On("AdCreated", inputAd).Once()
}

var (
	createExampleAd = types.Ad{
		ImageURL:        "https://via.placeholder.com/300x300",
		ClickThroughURL: "https://github.com",
	}

	createExamplePersistedAd = types.Ad{
		ID:              1,
		ImageURL:        "https://via.placeholder.com/300x300",
		ClickThroughURL: "https://github.com",
	}

	createExamplePayload = commands.CreatePayload{
		Ad: createExampleAd,
	}

	createExamplePersisterError = fmt.Errorf("Some error with the data store!")
)

type createFixtures struct {
	command   commands.Create
	persister *MockPersister
	notifier  *MockNotifier
}

func (f createFixtures) assertMockExpectations(t *testing.T) {
	f.persister.AssertExpectations(t)
	f.notifier.AssertExpectations(t)
}

func createSetup() createFixtures {
	p := new(MockPersister)
	n := new(MockNotifier)
	c := commands.Create{
		Persister: p,
		Notifier:  n,
	}

	return createFixtures{c, p, n}
}

func TestCreateOK(t *testing.T) {
	f := createSetup()
	f.persister.ExpectCreateSuccess(createExampleAd, createExamplePersistedAd)
	f.notifier.ExpectAdCreated(createExamplePersistedAd)

	err := f.command.Execute(createExamplePayload)
	assert.NoError(t, err)
	f.assertMockExpectations(t)
}

func TestCreatePersisterFailure(t *testing.T) {
	f := createSetup()
	f.persister.ExpectCreateError(createExampleAd, createExamplePersisterError)

	err := f.command.Execute(createExamplePayload)
	assert.True(t, errors.Is(err, createExamplePersisterError))
	f.assertMockExpectations(t)
}
