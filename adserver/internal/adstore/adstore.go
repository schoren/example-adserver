package adstore

import (
	"fmt"

	"github.com/schoren/example-adserver/pkg/types"
)

// Getter can fetch ads from the store
type Getter interface {
	Get(id int) (types.Ad, error)
}

// Setter can set new data into the store
type Setter interface {
	Set(types.Ad)
}

// GetSetter can get and set ads from the store
type GetSetter interface {
	Getter
	Setter
}

var ErrNotFound = fmt.Errorf("Ad not found")
