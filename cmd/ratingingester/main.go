package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"movieexample.com/rating/pkg/model"
)

// The main function creates a Kafka producer, reads rating events from a file, and sends them to a Kafka topic.
func main() {
	fmt.Println("Creating a Kafka producer")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
	}
	defer producer.Close()

	const fileName = "ratingsdata.json"
	log.Println("Reading data from file", fileName)

	ratingEvents, err := readRatingEvents(fileName)
	log.Printf("RatingData read from file: %v\n", ratingEvents)
	if err != nil {
		log.Fatalf("Failed to read rating events: %s\n", err)
	}

	const topic = "ratings"
	log.Println("Sending rating events to topic", topic)

	if err := produceRatingEvents(producer, topic, ratingEvents); err != nil {
		panic(err)
	}

	const timeout = 10 * time.Second
	log.Printf("Waiting for %s to send rating events\n", timeout)

	producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string) ([]model.RatingEvent, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ratingEvents []model.RatingEvent
	if err := json.NewDecoder(f).Decode(&ratingEvents); err != nil {
		return nil, err
	}

	return ratingEvents, nil
}

func produceRatingEvents(producer *kafka.Producer, topic string, ratingEvents []model.RatingEvent) error {
	for _, ratingEvent := range ratingEvents {
		encodedEvent, err := json.Marshal(ratingEvent)
		if err != nil {
			return err
		}

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          encodedEvent,
		}
		if err := producer.Produce(message, nil); err != nil {
			return err
		}
	}
	return nil
}
