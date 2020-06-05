package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/schoren/example-adserver/ads/internal/types"
)

// Notifier notifier
type Notifier struct {
	producer sarama.SyncProducer
}

func NewNotifier(producer sarama.SyncProducer) *Notifier {
	return &Notifier{producer}
}

type ad struct {
	ID              int    `json:"id"`
	ImageURL        string `json:"image_url"`
	ClickThroughURL string `json:"clickthrough_url"`
}

// AdUpdate notifies subscriber about changes in ads
func (n *Notifier) AdUpdate(inputAd types.Ad) {
	encoded, _ := json.Marshal(inputAd)
	_, _, err := n.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "ad-updates",
		Value: sarama.StringEncoder(string(encoded)),
	})

	if err != nil {
		log.Printf("Failed to produce kafka message: %s", err.Error())
	}
}
