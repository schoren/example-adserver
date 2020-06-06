package handlers_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/adserver/internal/handlers"
	"github.com/schoren/example-adserver/adserver/internal/renderer"
	"github.com/schoren/example-adserver/testutil/http/request"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	serveExampleRenderedAd = `<a href="http://example.org/"><img src="http://example.org/img.gif"></a>`

	serveExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org/",
	}
)

type StringRenderer struct {
	s string
}

func (s StringRenderer) Render() string {
	return s.s
}

type MockServe struct {
	mock.Mock
}

func (m *MockServe) Execute(id int) (renderer.Renderer, error) {
	args := m.Called(id)
	return args.Get(0).(renderer.Renderer), args.Error(1)
}

func serveSetup() (*mux.Router, *MockServe) {
	router := mux.NewRouter()
	serve := new(MockServe)

	handlers.ConfigureRouter(router)
	handlers.ServeCommand = serve

	return router, serve
}

func buildServeRequest(t *testing.T, id string) *http.Request {
	url := strings.Replace(handlers.ServeURL, "{id}", id, 1)
	r := request.MustBuild(t, handlers.ServeMethod, url, nil)
	return r
}

func TestServeOK(t *testing.T) {
	router, serve := serveSetup()
	serve.On("Execute", serveExampleAd.ID).Return(StringRenderer{serveExampleRenderedAd}, nil)

	req := buildServeRequest(t, strconv.Itoa(serveExampleAd.ID))
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, serveExampleRenderedAd, rr.Body.String())
}

func TestServeInvalidID(t *testing.T) {
	invalidIDs := []string{"0", "str"}
	for _, tt := range invalidIDs {
		t.Run(tt, func(t *testing.T) {
			router, serve := serveSetup()
			req := buildServeRequest(t, tt)
			rr := request.Exec(req, router.ServeHTTP)

			assert.Equal(t, http.StatusNotFound, rr.Code)
			assert.Empty(t, rr.Body.String())
			serve.AssertExpectations(t)
		})
	}
}

func TestServeCommandError(t *testing.T) {
	router, serve := serveSetup()
	serve.On("Execute", serveExampleAd.ID).Return(StringRenderer{}, fmt.Errorf("Some command error"))

	req := buildServeRequest(t, strconv.Itoa(serveExampleAd.ID))
	rr := request.Exec(req, router.ServeHTTP)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Empty(t, rr.Body.String())
	serve.AssertExpectations(t)
}