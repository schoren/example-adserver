package commands

import (
	"fmt"

	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/adserver/internal/renderer"
)

type emptyRenderer struct{}

func (r emptyRenderer) Render() string { return "" }

// Serve tries to create a Renderer from the given ad ID
type Serve struct {
	adStore adstore.Getter
}

func NewServe(adStore adstore.Getter) *Serve {
	return &Serve{adStore: adStore}
}

func (c *Serve) Execute(adID int) (renderer.Renderer, error) {
	ad, err := c.adStore.Get(adID)
	if err != nil {
		return emptyRenderer{}, fmt.Errorf("Cannot get ad with ID %d: %w", adID, err)
	}

	return renderer.Create(ad)
}
