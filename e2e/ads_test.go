// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	adServiceBaseURL string

	adCreateURL string
)

const (
	createAdJSONRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`

	exampleAdServed = `<a href="https://github.com"><img src="https://via.placeholder.com/300x300"></a>`
)

func init() {
	adServiceBaseURL = os.Getenv("AD_SERVICE_BASE_URL")

	adCreateURL = fmt.Sprintf("%s/", adServiceBaseURL)
}

func TestAdCreateIsServedCorrectly(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(createAdJSONRequest).
		Post(adCreateURL)
	if err != nil {
		t.Fatalf("Failed to create ad: %s", err.Error())
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode())

	adURL := resp.Header().Get("Location")
	assert.NotEmpty(t, adURL)

	// wait for the message to propagate
	time.Sleep(5 * time.Second)

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(createAdJSONRequest).
		Get(adURL)
	if err != nil {
		t.Fatalf("Failed to get ad: %s", err.Error())
	}

	assert.Equal(t, exampleAdServed, string(resp.Body()))
}
