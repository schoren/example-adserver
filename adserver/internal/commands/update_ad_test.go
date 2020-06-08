package commands_test

import (
	"testing"

	"github.com/schoren/example-adserver/adserver/internal/commands"
	"github.com/schoren/example-adserver/types"
	"github.com/stretchr/testify/assert"
)

var (
	updateAdExampleAd = types.Ad{
		ID:              1,
		ImageURL:        "http://example.org/img.gif",
		ClickThroughURL: "http://example.org",
	}

	updateAdExamplePayload = commands.UpdateAdPayload{
		Ad: updateAdExampleAd,
	}
)

func setupUpdateAd() (commands.UpdateAd, *MockAdStore) {
	mockAdStore := new(MockAdStore)
	cmd := commands.NewUpdateAd(mockAdStore)

	return cmd, mockAdStore
}

func TestUpdateAd(t *testing.T) {
	cmd, as := setupUpdateAd()
	as.ExpectSet(updateAdExampleAd)

	err := cmd.Execute(updateAdExamplePayload)

	assert.NoError(t, err)
	as.AssertExpectations(t)
}
