package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/pkg/testutil/http/request"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockActiveLister struct {
	mock.Mock
}

func (m *MockActiveLister) ListActive() ([]types.Ad, error) {
	args := m.Called()
	return args.Get(0).([]types.Ad), args.Error(1)
}

func (m *MockActiveLister) ExpectListActiveSuccess(ads []types.Ad) {
	m.On("ListActive").Return(ads, nil)
}

func (m *MockActiveLister) ExpectListActiveError(err error) {
	m.On("ListActive").Return([]types.Ad{}, err)
}

func buildListActiveRequest(t *testing.T) *http.Request {
	r := request.MustBuild(t, handlers.ListActiveMethod, handlers.ListActiveURL, nil)
	return r
}

func listActiveSetup() (*mux.Router, *MockActiveLister) {
	router := mux.NewRouter()
	action := new(MockActiveLister)
	h := handlers.NewListActive(action)
	h.Register(router)

	return router, action
}

func TestListActive(t *testing.T) {
	t.Parallel()

	var (
		exampleAds = []types.Ad{
			{
				ID:              1,
				ImageURL:        "https://example.org/1.png",
				ClickThroughURL: "https://example.org/1.html",
			},
			{
				ID:              2,
				ImageURL:        "https://example.org/2.png",
				ClickThroughURL: "https://example.org/2.html",
			},
		}

		exampleAdsJSON = `[{"ID":1,"ImageURL":"https://example.org/1.png","ClickThroughURL":"https://example.org/1.html"},{"ID":2,"ImageURL":"https://example.org/2.png","ClickThroughURL":"https://example.org/2.html"}]`

		exampleError = fmt.Errorf("There was some error")
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		router, lister := listActiveSetup()
		lister.ExpectListActiveSuccess(exampleAds)

		req := buildListActiveRequest(t)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, exampleAdsJSON, rr.Body.String())
		lister.AssertExpectations(t)
	})

	t.Run("Action Error", func(t *testing.T) {
		t.Parallel()

		router, lister := listActiveSetup()
		lister.ExpectListActiveError(exampleError)

		req := buildListActiveRequest(t)
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Empty(t, rr.Body.String())
		lister.AssertExpectations(t)
	})
}
