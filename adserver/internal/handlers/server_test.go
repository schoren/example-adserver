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
	"github.com/schoren/example-adserver/pkg/testutil/http/request"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockServe) Serve(id int) (renderer.Renderer, error) {
	args := m.Called(id)
	return args.Get(0).(renderer.Renderer), args.Error(1)
}

func serveSetup() (*mux.Router, *MockServe) {
	router := mux.NewRouter()
	action := new(MockServe)
	h := handlers.NewServer(action)
	h.Register(router)

	return router, action
}

func buildServeRequest(t *testing.T, id string) *http.Request {
	url := strings.Replace(handlers.ServeURL, "{id}", id, 1)
	r := request.MustBuild(t, handlers.ServeMethod, url, nil)
	return r
}

func TestServe(t *testing.T) {
	t.Parallel()

	var (
		exampleRenderedAd = `<a href="http://example.org/"><img src="http://example.org/img.gif"></a>`

		exampleAd = types.Ad{
			ID:              1,
			ImageURL:        "http://example.org/img.gif",
			ClickThroughURL: "http://example.org/",
		}
	)

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		router, serve := serveSetup()
		serve.On("Serve", exampleAd.ID).Return(StringRenderer{exampleRenderedAd}, nil)

		req := buildServeRequest(t, strconv.Itoa(exampleAd.ID))
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, exampleRenderedAd, rr.Body.String())
	})

	t.Run("Action error", func(t *testing.T) {
		t.Parallel()

		router, serve := serveSetup()
		serve.On("Serve", exampleAd.ID).Return(StringRenderer{}, fmt.Errorf("Some command error"))

		req := buildServeRequest(t, strconv.Itoa(exampleAd.ID))
		rr := request.Exec(req, router.ServeHTTP)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Empty(t, rr.Body.String())
		serve.AssertExpectations(t)
	})

	t.Run("Invalid IDs", func(t *testing.T) {
		t.Parallel()

		invalidIDs := []string{"0", "str"}
		for _, tt := range invalidIDs {
			t.Run(tt, func(t *testing.T) {
				t.Parallel()
				router, serve := serveSetup()
				req := buildServeRequest(t, tt)
				rr := request.Exec(req, router.ServeHTTP)

				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.Empty(t, rr.Body.String())
				serve.AssertExpectations(t)
			})
		}
	})
}
