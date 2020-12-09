package actions

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/types"
)

type CreatePersister interface {
	Create(types.Ad) (types.Ad, error)
}

type Creator interface {
	Create(ad types.Ad) (types.Ad, error)
}

func NewCreator(p CreatePersister, n notifier, i instrumentation.Instrumentator) Creator {
	return &create{
		persister:      p,
		notifier:       n,
		instrumentator: i,
	}
}

type create struct {
	persister      CreatePersister
	notifier       notifier
	instrumentator instrumentation.Instrumentator
}

func (a *create) Create(ad types.Ad) (types.Ad, error) {
	a.instrumentator.OnStart()
	defer a.instrumentator.OnComplete()

	ad, err := a.persister.Create(ad)
	if err != nil {
		err = fmt.Errorf("error creating ad: %w", err)
		a.instrumentator.OnError(err)
		return types.Ad{}, err
	}

	a.notifier.AdUpdate(ad)

	return ad, nil
}
