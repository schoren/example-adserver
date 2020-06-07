package renderer_test

import (
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/renderer"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
)

var (
	rendererCreateExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org/",
	}

	rendererCreateExampleEmptyAd = types.Ad{}

	rendererCreateExampleRenderedAd = `<a href="http://example.org/"><img src="http://example.org/img.gif"></a>`
)

func TestCreate(t *testing.T) {
	r, err := renderer.Create(rendererCreateExampleAd)

	assert.NoError(t, err)
	assert.Equal(t, rendererCreateExampleRenderedAd, r.Render())
}
