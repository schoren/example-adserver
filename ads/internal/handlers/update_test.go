package handlers_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/schoren/example-adserver/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/pkg/testutil/http/request"
)

var (
	updateExampleID = "1"

	updateExampleInvalidJSONRequest = `{invalid}`

	updateExampleOKRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`

	updateExampleCommandPayload = commands.UpdatePayload{
		Ad: types.Ad{
			ID:              1,
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		},
	}

	updateExampleCommandError = fmt.Errorf("There was some error")
)

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) Execute(p commands.UpdatePayload) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockUpdater) ExpectExecuteSuccess(payload commands.UpdatePayload) {
	m.On("Execute", updateExampleCommandPayload).Return(nil)
}

func (m *MockUpdater) ExpectExecuteError(payload commands.UpdatePayload) {
	m.On("Execute", updateExampleCommandPayload).Return(updateExampleCommandError)
}

func updateSetup() (*mux.Router, *MockUpdater) {
	router := mux.NewRouter()
	handlers.ConfigureRouter(router)
	updater := new(MockUpdater)
	handlers.UpdateCommand = updater

	return router, updater
}

func buildUpdateRequest(t *testing.T, id, body string) *http.Request {
	url := strings.Replace(handlers.UpdateURL, "{id}", id, 1)
	r := request.MustBuild(t, handlers.UpdateMethod, url, request.Body(body))
	request.IsJSON(r)
	return r
}

func TestUpdateSuccess(t *testing.T) {
	router, updater := updateSetup()
	updater.ExpectExecuteSuccess(updateExampleCommandPayload)

	req := buildUpdateRequest(t, updateExampleID, updateExampleOKRequest)
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
	updater.AssertExpectations(t)
}

func TestUpdateInvalidJSON(t *testing.T) {
	router, updater := updateSetup()

	req := buildUpdateRequest(t, updateExampleID, updateExampleInvalidJSONRequest)
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.NotEmpty(t, rr.Body.String())
	updater.AssertExpectations(t)
}

func TestUpdateCommandError(t *testing.T) {
	router, updater := updateSetup()
	updater.ExpectExecuteError(updateExampleCommandPayload)

	req := buildUpdateRequest(t, updateExampleID, updateExampleOKRequest)
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Empty(t, rr.Body.String())
	updater.AssertExpectations(t)
}

func TestUpdateInvalidID(t *testing.T) {
	invalidIDs := []string{"0", "str"}
	for _, tt := range invalidIDs {
		t.Run(tt, func(t *testing.T) {
			router, updater := updateSetup()
			req := buildUpdateRequest(t, tt, updateExampleOKRequest)
			rr := request.Exec(req, router.ServeHTTP)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.NotEmpty(t, rr.Body.String())
			updater.AssertExpectations(t)
		})
	}
}
