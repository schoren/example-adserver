// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	adServiceBaseURL string
)

const (
	createAdJSONRequest = `{
		"image_url":"https://via.placeholder.com/300x300",
		"clickthrough_url":"https://github.com"
	}`
)

func init() {
	adServiceBaseURL = os.Getenv("AD_SERVICE_BASE_URL")
}

func TestAdCreateIsServedCorrectly(t *testing.T) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(createAdJSONRequest).
		Post(fmt.Sprintf("%s/", adServiceBaseURL))
	if err != nil {
		t.Fatalf("Cannot make request: %s", err.Error())
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode())
}
