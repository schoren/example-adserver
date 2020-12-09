package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/pkg/testutil/http/request"
	"github.com/schoren/example-adserver/pkg/types"
)

type MockCreater struct {
	mock.Mock
}

func (m *MockCreater) Create(ad types.Ad) (types.Ad, error) {
	args := m.Called(ad)
	return args.Get(0).(types.Ad), args.Error(1)
}

func (m *MockCreater) ExpectCreateSuccess(ad, newAd types.Ad) {
	m.On("Create", ad).Return(newAd, nil)
}

func (m *MockCreater) ExpectCreateError(err error) {
	m.On("Create", mock.Anything).Return(types.Ad{}, err)
}

func createSetup(baseURL string) (*mux.Router, *MockCreater) {
	router := mux.NewRouter()
	action := new(MockCreater)
	h := handlers.NewCreate(action, baseURL)
	h.Register(router)

	return router, action
}

func buildCreateRequest(t *testing.T, body string) *http.Request {
	r := request.MustBuild(t, handlers.CreateMethod, handlers.CreateURL, request.Body(body))
	request.IsJSON(r)
	return r
}

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		exampleInvalidJSONRequest = `{invalid}`

		exampleOKRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`

		exampleAd = types.Ad{
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		exampleNewAd = types.Ad{
			ID:              1,
			ImageURL:        "https://via.placeholder.com/300x300",
			ClickThroughURL: "https://github.com",
		}

		exampleError   = fmt.Errorf("There was some error")
		exampleBaseURL = "http://localhost"
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		router, creater := createSetup(exampleBaseURL)
		creater.ExpectCreateSuccess(exampleAd, exampleNewAd)

		req := buildCreateRequest(t, exampleOKRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Empty(t, rr.Body.String())
		assert.Equal(t, fmt.Sprintf("%s/%d", exampleBaseURL, exampleNewAd.ID), rr.Header().Get("Location"))
		creater.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		t.Parallel()

		router, creater := createSetup(exampleBaseURL)

		req := buildCreateRequest(t, exampleInvalidJSONRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
		creater.AssertExpectations(t)
	})

	t.Run("Action error", func(t *testing.T) {
		t.Parallel()

		router, creater := createSetup(exampleBaseURL)
		creater.ExpectCreateError(exampleError)

		req := buildCreateRequest(t, exampleOKRequest)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, rr.Body.String())
		creater.AssertExpectations(t)
	})
}
