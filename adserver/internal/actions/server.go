package actions

import (
	"fmt"

	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/adserver/internal/renderer"
)

type Server interface {
	Serve(adID int) (renderer.Renderer, error)
}

func NewServer(adStore adstore.Getter) Server {
	return &server{adStore: adStore}
}

type server struct {
	adStore adstore.Getter
}

func (c *server) Serve(adID int) (renderer.Renderer, error) {
	ad, err := c.adStore.Get(adID)
	if err != nil {
		return emptyRenderer{}, fmt.Errorf("Cannot get ad with ID %d: %w", adID, err)
	}

	return renderer.Create(ad)
}

type emptyRenderer struct{}

func (r emptyRenderer) Render() string { return "" }
