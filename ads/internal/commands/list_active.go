package commands

import "github.com/schoren/example-adserver/types"

type ActiveAdGetter interface {
	GetActive() ([]types.Ad, error)
}

type ListActiveCommand struct {
	adGetter ActiveAdGetter
}

func NewListActiveCommand(adGetter ActiveAdGetter) *ListActiveCommand {
	return &ListActiveCommand{
		adGetter: adGetter,
	}
}

func (c *ListActiveCommand) Execute() ([]types.Ad, error) {
	return c.adGetter.GetActive()
}
