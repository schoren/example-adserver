package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/schoren/example-adserver/adserver/internal/actions"
	"github.com/schoren/example-adserver/pkg/types"
)

func NewAdUpdater(updater actions.AdUpdater) *AdUpdater {
	return &AdUpdater{
		Ready:   make(chan bool),
		updater: updater,
	}
}

// AdUpdater represents a Sarama consumer group consumer
type AdUpdater struct {
	Ready   chan bool
	updater actions.AdUpdater
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (adupdater *AdUpdater) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(adupdater.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (adupdater *AdUpdater) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (adupdater *AdUpdater) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var updatedAd types.Ad
		err := json.Unmarshal(message.Value, &updatedAd)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
		}
		adupdater.updater.Update(updatedAd)

		session.MarkMessage(message, "")
	}

	return nil
}
