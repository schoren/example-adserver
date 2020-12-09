package actions

import "github.com/schoren/example-adserver/pkg/types"

// notifier can propagate events to other components of the system
type notifier interface {
	AdUpdate(types.Ad)
}
