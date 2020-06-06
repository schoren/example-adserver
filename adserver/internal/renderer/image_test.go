package renderer_test

import (
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/renderer"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
)

var (
	imageRendererExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org/",
	}

	imageRendererExampleEmptyAd = types.Ad{}

	imageRendererExampleRenderedAd = `<a href="http://example.org/"><img src="http://example.org/img.gif"></a>`
)

func TestImageRenderer(t *testing.T) {
	r, err := renderer.NewImage(imageRendererExampleAd)

	assert.NoError(t, err)
	assert.Equal(t, imageRendererExampleRenderedAd, r.Render())
}

func TestImageRendererAdValidation(t *testing.T) {
	r, err := renderer.NewImage(imageRendererExampleEmptyAd)

	assert.Error(t, err)
	assert.Empty(t, r)
}
