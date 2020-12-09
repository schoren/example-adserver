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

var (
	listActiveExampleAds = []types.Ad{
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

	listActiveExampleAdsJSON = `[{"ID":1,"ImageURL":"https://example.org/1.png","ClickThroughURL":"https://example.org/1.html"},{"ID":2,"ImageURL":"https://example.org/2.png","ClickThroughURL":"https://example.org/2.html"}]`

	listActiveExampleCommandError = fmt.Errorf("There was some error")
)

type MockActiveLister struct {
	mock.Mock
}

func (m *MockActiveLister) Execute() ([]types.Ad, error) {
	args := m.Called()
	return args.Get(0).([]types.Ad), args.Error(1)
}

func (m *MockActiveLister) ExpectExecuteSuccess() {
	m.On("Execute").Return(listActiveExampleAds, nil)
}

func (m *MockActiveLister) ExpectExecuteError() {
	m.On("Execute").Return([]types.Ad{}, listActiveExampleCommandError)
}

func buildListActiveRequest(t *testing.T) *http.Request {
	r := request.MustBuild(t, handlers.ListActiveMethod, handlers.ListActiveURL, nil)
	return r
}

func listActiveSetup() (*mux.Router, *MockActiveLister) {
	router := mux.NewRouter()
	handlers.ConfigureRouter(router)

	lister := new(MockActiveLister)
	handlers.ListActiveCommand = lister

	return router, lister
}

func TestListActiveSuccess(t *testing.T) {
	router, lister := listActiveSetup()
	lister.ExpectExecuteSuccess()

	req := buildListActiveRequest(t)
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, listActiveExampleAdsJSON, rr.Body.String())
	lister.AssertExpectations(t)
}

func TestListActiveListerError(t *testing.T) {
	router, lister := listActiveSetup()
	lister.ExpectExecuteError()

	req := buildListActiveRequest(t)
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Empty(t, rr.Body.String())
	lister.AssertExpectations(t)
}
