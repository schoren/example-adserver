package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/schoren/example/ads/ads/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example/ads/ads/internal/commands"
	"github.com/schoren/example/ads/ads/internal/handlers"
	"github.com/schoren/example/ads/testutil/http/request"
)

var (
	createExampleInvalidJSONRequest = `{invalid}`

	createExampleOKRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`

	createExampleCommandPayload = commands.CreatePayload{
		Ad: types.Ad{
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		},
	}

	createExampleCommandError = fmt.Errorf("There was some error")
)

type MockCreater struct {
	mock.Mock
}

func (m *MockCreater) Execute(p commands.CreatePayload) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockCreater) ExpectExecuteSuccess(payload commands.CreatePayload) {
	m.On("Execute", createExampleCommandPayload).Return(nil)
}

func (m *MockCreater) ExpectExecuteError(payload commands.CreatePayload) {
	m.On("Execute", createExampleCommandPayload).Return(createExampleCommandError)
}

func createSetup() *MockCreater {
	creater := new(MockCreater)
	handlers.CreateCommand = creater

	return creater
}

func buildCreateRequest(t *testing.T, body string) *http.Request {
	r := request.MustBuild(t, handlers.CreateMethod, handlers.CreateURL, request.Body(body))
	request.IsJSON(r)
	return r
}

func TestCreateSuccess(t *testing.T) {
	creater := createSetup()
	creater.ExpectExecuteSuccess(createExampleCommandPayload)

	req := buildCreateRequest(t, createExampleOKRequest)
	rr := request.Exec(req, http.HandlerFunc(handlers.Create))

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Empty(t, rr.Body.String())
	creater.AssertExpectations(t)
}

func TestCreateInvalidJSON(t *testing.T) {
	creater := createSetup()

	req := buildCreateRequest(t, createExampleInvalidJSONRequest)
	rr := request.Exec(req, http.HandlerFunc(handlers.Create))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	creater.AssertExpectations(t)
}

func TestCreateCommandError(t *testing.T) {
	creater := createSetup()
	creater.ExpectExecuteError(createExampleCommandPayload)

	req := buildCreateRequest(t, createExampleOKRequest)
	rr := request.Exec(req, http.HandlerFunc(handlers.Create))

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Empty(t, rr.Body.String())
	creater.AssertExpectations(t)
}
