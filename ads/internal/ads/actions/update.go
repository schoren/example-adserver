package actions

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/types"
)

type UpdatePersister interface {
	Update(types.Ad) error
}

type Updater interface {
	Update(ad types.Ad) error
}

func NewUpdater(p UpdatePersister, n notifier, i instrumentation.Instrumentator) Updater {
	return &update{
		updater:        p,
		notifier:       n,
		instrumentator: i,
	}
}

type update struct {
	updater        UpdatePersister
	notifier       notifier
	instrumentator instrumentation.Instrumentator
}

func (a *update) Update(ad types.Ad) error {
	a.instrumentator.OnStart()
	defer a.instrumentator.OnComplete()

	err := a.updater.Update(ad)
	if err != nil {
		err = fmt.Errorf("error updating ad: %w", err)
		a.instrumentator.OnError(err)
		return err
	}

	a.notifier.AdUpdate(ad)

	return nil
}
