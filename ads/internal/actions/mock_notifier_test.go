package actions_test

import (
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/mock"
)

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) AdUpdate(inputAd types.Ad) {
	m.Called(inputAd)
}

func (m *MockNotifier) ExpectAdUpdate(inputAd types.Ad) {
	m.On("AdUpdate", inputAd).Once()
}
