package actions_test

import (
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/actions"
	"github.com/schoren/example-adserver/pkg/types"
	"github.com/stretchr/testify/assert"
)

func setupUpdateAd() (actions.AdUpdater, *MockAdStore) {
	mockAdStore := new(MockAdStore)
	action := actions.NewAdUpdater(mockAdStore)

	return action, mockAdStore
}

func TestAdUpdater(t *testing.T) {
	t.Parallel()

	var (
		exampleAd = types.Ad{
			ID:              1,
			ImageURL:        "http://example.org/img.gif",
			ClickThroughURL: "http://example.org",
		}
	)

	action, as := setupUpdateAd()
	as.ExpectSet(exampleAd)

	err := action.Update(exampleAd)

	assert.NoError(t, err)
	as.AssertExpectations(t)
}
