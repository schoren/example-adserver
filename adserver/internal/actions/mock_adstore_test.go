package actions_test

import (
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/mock"
)

type MockAdStore struct {
	mock.Mock
}

func (m *MockAdStore) Get(id int) (types.Ad, error) {
	args := m.Called(id)
	return args.Get(0).(types.Ad), args.Error(1)
}

func (m *MockAdStore) Set(ad types.Ad) {
	m.Called(ad)
}

func (m *MockAdStore) ExpectGetSuccess(id int, ad types.Ad) {
	m.On("Get", id).Return(ad, nil)
}

func (m *MockAdStore) ExpectGetError(id int) {
	m.On("Get", id).Return(types.Ad{}, serveExampleAdStoreError)
}

func (m *MockAdStore) ExpectSet(ad types.Ad) {
	m.On("Set", ad)
}
