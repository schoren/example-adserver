package renderer

import "github.com/schoren/example-adserver/pkg/types"

// Renderer can return a string representation of an ad, ready to be served
type Renderer interface {
	Render() string
}

// Create returns a new renderer based on the details of the ad
func Create(ad types.Ad) (Renderer, error) {
	return NewImage(ad)
}
