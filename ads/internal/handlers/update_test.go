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

	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/pkg/testutil/http/request"
)

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) Update(ad types.Ad) error {
	args := m.Called(ad)
	return args.Error(0)
}

func (m *MockUpdater) ExpectUpdateSuccess(ad types.Ad) {
	m.On("Update", ad).Return(nil)
}

func (m *MockUpdater) ExpectUpdateError(err error) {
	m.On("Update", mock.Anything).Return(err)
}

func updateSetup() (*mux.Router, *MockUpdater) {
	router := mux.NewRouter()
	action := new(MockUpdater)
	h := handlers.NewUpdate(action)
	h.Register(router)

	return router, action
}

func buildUpdateRequest(t *testing.T, id, body string) *http.Request {
	url := strings.Replace(handlers.UpdateURL, "{id}", id, 1)
	r := request.MustBuild(t, handlers.UpdateMethod, url, request.Body(body))
	request.IsJSON(r)
	return r
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	var (
		exampleID = "1"

		exampleInvalidJSONRequest = `{invalid}`

		exampleOKRequest = `{
			"image_url":"https://via.placeholder.com/300x300",
			"clickthrough_url":"https://github.com"
		}`

		exampleAd = types.Ad{
			ID:              1,
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		exampleError = fmt.Errorf("There was some error")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		router, updater := updateSetup()
		updater.ExpectUpdateSuccess(exampleAd)

		req := buildUpdateRequest(t, exampleID, exampleOKRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
		updater.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		router, updater := updateSetup()

		req := buildUpdateRequest(t, exampleID, exampleInvalidJSONRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
		updater.AssertExpectations(t)
	})

	t.Run("Action Error", func(t *testing.T) {
		t.Parallel()

		router, updater := updateSetup()
		updater.ExpectUpdateError(exampleError)

		req := buildUpdateRequest(t, exampleID, exampleOKRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, rr.Body.String())
		updater.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		t.Parallel()

		invalidIDs := []string{"0", "str"}
		for _, tt := range invalidIDs {
			t.Run(tt, func(t *testing.T) {
				t.Parallel()

				router, updater := updateSetup()
				req := buildUpdateRequest(t, tt, exampleOKRequest)
				rr := request.Exec(req, router.ServeHTTP)

				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.NotEmpty(t, rr.Body.String())
				updater.AssertExpectations(t)
			})
		}
	})
}
