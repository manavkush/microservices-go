package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"movieexample.com/rating/pkg/model"
)

// Ingester defines a Kafka ingester
type Ingester struct {
	consumer *kafka.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester
func NewIngester(adrr string, groupID string, topic string) (*Ingester, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": adrr,
		"group.id":          groupID,
	})
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer, topic}, nil
}

// Ingest starts ingesting from Kafka and returns a channel of rating events
// representing the data consumed from topics
func (i *Ingester) Ingest(ctx context.Context) (<-chan model.RatingEvent, error) {
	if err := i.consumer.SubscribeTopics([]string{i.topic}, nil); err != nil {
		return nil, err
	}

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
			default:
			}

			// Waiting indefinitely
			msg, err := i.consumer.ReadMessage(-1)
			if err != nil {
				log.Println("Consumer error:", err)
				continue
			}
			var event model.RatingEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("Failed to unmarshal event:", err)
				continue
			}
			ch <- event
		}
	}()

	return ch, nil
}
