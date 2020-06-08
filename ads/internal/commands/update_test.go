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

func (m *MockUpdatePersister) ExpectUpdateError(inputAd types.Ad, err error) {
	m.On("Update", inputAd).Return(err)
}

var (
	updateExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "https://via.placeholder.com/300x300",
		ClickThroughURL: "https://github.com",
	}

	updateExamplePayload = commands.UpdatePayload{
		Ad: updateExampleAd,
	}

	updateExamplePersisterError = fmt.Errorf("Some error with the data store!")
)

type UpdateFixtures struct {
	command   *commands.Update
	persister *MockUpdatePersister
	notifier  *MockNotifier
}

func (f UpdateFixtures) assertMockExpectations(t *testing.T) {
	f.persister.AssertExpectations(t)
	f.notifier.AssertExpectations(t)
}

func UpdateSetup() UpdateFixtures {
	p := new(MockUpdatePersister)
	n := new(MockNotifier)
	c := commands.NewUpdate(p, n)

	return UpdateFixtures{c, p, n}
}

func TestUpdateOK(t *testing.T) {
	f := UpdateSetup()
	f.persister.ExpectUpdateSuccess(updateExampleAd)
	f.notifier.ExpectAdUpdate(updateExampleAd)

	err := f.command.Execute(updateExamplePayload)

	assert.NoError(t, err)
	f.assertMockExpectations(t)
}

func TestUpdatePersisterFailure(t *testing.T) {
	f := UpdateSetup()
	f.persister.ExpectUpdateError(updateExampleAd, updateExamplePersisterError)

	err := f.command.Execute(updateExamplePayload)

	assert.True(t, errors.Is(err, updateExamplePersisterError))
	f.assertMockExpectations(t)
}
