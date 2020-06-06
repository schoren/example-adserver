package commands

import (
	"fmt"

	"github.com/schoren/example-adserver/adserver/internal/renderer"
	"github.com/schoren/example-adserver/types"
)

type emptyRenderer struct{}

func (r emptyRenderer) Render() string { return "" }

type AdStoreGetter interface {
	Get(id int) (types.Ad, error)
}

// ServeCommand tries to create a Renderer from the given ad ID
type ServeCommand struct {
	AdStore AdStoreGetter
}

func (c *ServeCommand) Execute(adID int) (renderer.Renderer, error) {
	ad, err := c.AdStore.Get(adID)
	if err != nil {
		return emptyRenderer{}, fmt.Errorf("Cannot get ad with ID %d: %w", adID, err)
	}

	return renderer.Create(ad)
}
