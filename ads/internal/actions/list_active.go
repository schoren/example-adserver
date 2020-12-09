package actions

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/types"
)

type ActiveAdGetter interface {
	GetActive() ([]types.Ad, error)
}

type ActiveLister interface {
	ListActive() ([]types.Ad, error)
}

func NewActiveLister(ag ActiveAdGetter, i instrumentation.Instrumentator) ActiveLister {
	return &listActive{
		adGetter:       ag,
		instrumentator: i,
	}
}

type listActive struct {
	adGetter       ActiveAdGetter
	instrumentator instrumentation.Instrumentator
}

func (a *listActive) ListActive() ([]types.Ad, error) {
	a.instrumentator.OnStart()
	defer a.instrumentator.OnComplete()

	ads, err := a.adGetter.GetActive()
	if err != nil {
		err = fmt.Errorf("error listing active ads: %w", err)
		a.instrumentator.OnError(err)
		return []types.Ad{}, err
	}

	return ads, nil
}
