package request

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Body(body string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(body))
}

func MustBuild(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func IsJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func Exec(req *http.Request, h http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
