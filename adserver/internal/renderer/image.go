package renderer

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/types"
)

// Image is a renderer that renders an Image ad as an HTML image
type Image struct {
	ad types.Ad
}

func adValidForImage(ad types.Ad) error {
	if ad.ImageURL == "" {
		return fmt.Errorf("ImageURL must be set")
	}

	if ad.ClickThroughURL == "" {
		return fmt.Errorf("ClickThroughURL must be set")
	}

	return nil
}

func NewImage(ad types.Ad) (Image, error) {
	if err := adValidForImage(ad); err != nil {
		return Image{}, fmt.Errorf("Ad is not renderable with Image: %w", err)
	}
	return Image{ad}, nil
}

func (r Image) Render() string {
	return fmt.Sprintf(`<a href="%s"><img src="%s"></a>`, r.ad.ClickThroughURL, r.ad.ImageURL)
}
