package renderer_test

import (
	"testing"

	"github.com/schoren/example-adserver/adserver/renderer"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
)

var (
	exampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org/",
	}

	exampleEmptyAd = types.Ad{}

	exampleRenderedAd = `<a href="http://example.org/"><img src="http://example.org/img.gif"></a>`
)

func TestImageRenderer(t *testing.T) {
	r, err := renderer.NewImage(exampleAd)

	assert.NoError(t, err)
	assert.Equal(t, exampleRenderedAd, r.Render())
}

func TestImageRendererAdValidation(t *testing.T) {
	r, err := renderer.NewImage(exampleEmptyAd)

	assert.Error(t, err)
	assert.Empty(t, r)
}
