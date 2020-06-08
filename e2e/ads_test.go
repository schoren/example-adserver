// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	adServiceBaseURL string
	adserverBaseURL  string

	adCreateURL string
	adUpdateURL string

	defaultExampleAdURL = ""
)

const (
	createAdJSONRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`

	updateAdJSONRequest = `{
		"image_url":"https://via.placeholder.com/100x100",
		"clickthrough_url":"https://github.com/?updated=1"
	}`

	exampleCreatedAdServed = `<a href="https://github.com"><img src="https://via.placeholder.com/300x300"></a>`
	exampleUpdatedAdServed = `<a href="https://github.com/?updated=1"><img src="https://via.placeholder.com/100x100"></a>`

	defaultAdExampleServed = `<a href="http://example.org/1.html"><img src="http://example.org/1.png"></a>`
)

func init() {
	adServiceBaseURL = os.Getenv("AD_SERVICE_BASE_URL")
	adserverBaseURL = os.Getenv("ADSERVER_BASE_URL")

	adCreateURL = fmt.Sprintf("%s/", adServiceBaseURL)
	adUpdateURL = fmt.Sprintf("%s/{id}", adServiceBaseURL)

	defaultExampleAdURL = fmt.Sprintf("%s/1", adserverBaseURL)
}

func getAdUpdateURL(adID string) string {
	return strings.Replace(adUpdateURL, "{id}", adID, 1)
}

func TestAdserverIsWarmedUpOnStart(t *testing.T) {
	t.Parallel()

	// We know that the DB is created with an example ad.
	// Since it's not created throught the API, it won't ever be propagated through kafka.
	// So, if the ad server has the ad, it was properly warmed up.

	client := resty.New()
	resp, err := client.R().Get(defaultExampleAdURL)
	if err != nil {
		t.Fatalf("Failed to get ad: %s", err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, defaultAdExampleServed, string(resp.Body()))

}

func TestAdCreateIsServedCorrectly(t *testing.T) {
	t.Parallel()

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

	assert.Equal(t, exampleCreatedAdServed, string(resp.Body()))
}

func TestAdUpdateIsServedCorrectly(t *testing.T) {
	t.Parallel()

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(createAdJSONRequest).
		Post(adCreateURL)
	if err != nil {
		t.Fatalf("Failed to create ad: %s", err.Error())
	}
	adURL := resp.Header().Get("Location")
	adID := strings.Replace(adURL, adserverBaseURL+"/", "", 1)

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(updateAdJSONRequest).
		Put(getAdUpdateURL(adID))
	if err != nil {
		t.Fatalf("Failed to update ad: %s", err.Error())
	}

	assert.Equal(t, http.StatusNoContent, resp.StatusCode())

	// wait for the message to propagate
	time.Sleep(5 * time.Second)

	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(createAdJSONRequest).
		Get(adURL)
	if err != nil {
		t.Fatalf("Failed to get ad: %s", err.Error())
	}

	assert.Equal(t, exampleUpdatedAdServed, string(resp.Body()))
}
