package commands

import "github.com/schoren/example-adserver/types"

type ActiveAdGetter interface {
	GetActive() ([]types.Ad, error)
}

type ListActive struct {
	adGetter ActiveAdGetter
}

func NewListActive(adGetter ActiveAdGetter) *ListActive {
	return &ListActive{
		adGetter: adGetter,
	}
}

func (c *ListActive) Execute() ([]types.Ad, error) {
	return c.adGetter.GetActive()
}
