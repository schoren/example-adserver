package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/schoren/example-adserver/types"
)

// Notifier notifier
type Notifier struct {
	producer sarama.SyncProducer
	topic    string
}

func NewNotifier(producer sarama.SyncProducer, topic string) *Notifier {
	return &Notifier{producer, topic}
}

// AdUpdate notifies subscriber about changes in ads
func (n *Notifier) AdUpdate(inputAd types.Ad) {
	encoded, _ := json.Marshal(inputAd)
	_, _, err := n.producer.SendMessage(&sarama.ProducerMessage{
		Topic: n.topic,
		Value: sarama.StringEncoder(string(encoded)),
	})

	if err != nil {
		log.Printf("Failed to produce kafka message: %s", err.Error())
	}
}
