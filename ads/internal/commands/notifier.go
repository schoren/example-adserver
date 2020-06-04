package commands

import "github.com/schoren/example/ads/ads/internal/types"

// Notifier can propagate events to other components of the system
type Notifier interface {
	AdUpdate(types.Ad)
}
