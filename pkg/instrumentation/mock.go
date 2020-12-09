package instrumentation

import (
	"errors"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) OnStart() {
	m.Called()
}

func (m *Mock) OnError(err error) {
	m.Called(err)
}

func (m *Mock) OnComplete() {
	m.Called()
}

func (m *Mock) ExpectOnStart() {
	m.On("OnStart")
}

func (m *Mock) ExpectOnError(err error) {
	m.On("OnError", mock.MatchedBy(func(inErr error) bool {
		return errors.Is(inErr, err)
	}))
}

func (m *Mock) ExpectOnComplete() {
	m.On("OnComplete")
}
