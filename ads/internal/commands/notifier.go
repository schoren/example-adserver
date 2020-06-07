package commands

import "github.com/schoren/example-adserver/types"

// Notifier can propagate events to other components of the system
type Notifier interface {
	AdUpdate(types.Ad)
}
