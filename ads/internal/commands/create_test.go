package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/types"
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
	command   *commands.Create
	persister *MockCreatePersister
	notifier  *MockNotifier
}

func (f createFixtures) assertMockExpectations(t *testing.T) {
	f.persister.AssertExpectations(t)
	f.notifier.AssertExpectations(t)
}

func createSetup() createFixtures {
	p := new(MockCreatePersister)
	n := new(MockNotifier)
	c := commands.NewCreate(p, n)

	return createFixtures{c, p, n}
}

func TestCreateOK(t *testing.T) {
	f := createSetup()
	f.persister.ExpectCreateSuccess(createExampleAd, createExamplePersistedAd)
	f.notifier.ExpectAdUpdate(createExamplePersistedAd)

	ad, err := f.command.Execute(createExamplePayload)
	assert.NoError(t, err)
	assert.Equal(t, createExamplePersistedAd, ad)
	f.assertMockExpectations(t)
}

func TestCreatePersisterFailure(t *testing.T) {
	f := createSetup()
	f.persister.ExpectCreateError(createExampleAd, createExamplePersisterError)

	ad, err := f.command.Execute(createExamplePayload)
	assert.True(t, errors.Is(err, createExamplePersisterError))
	assert.Equal(t, types.Ad{}, ad)
	f.assertMockExpectations(t)
}
