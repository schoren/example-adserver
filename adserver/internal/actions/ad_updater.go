package actions

import (
	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/pkg/types"
)

type AdUpdater interface {
	Update(types.Ad) error
}

func NewAdUpdater(adStore adstore.Setter) *adUpdater {
	return &adUpdater{adStore: adStore}
}

type adUpdater struct {
	adStore adstore.Setter
}

func (a *adUpdater) Update(ad types.Ad) error {
	a.adStore.Set(ad)
	return nil
}
