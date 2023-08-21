package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := []string{"https://psrc-10wzj.ap-southeast-2.aws.confluent.cloud"} // Replace with your broker addresses
	topic := "test-topic"                                                        // Replace with your desired topic

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 0,    // Set to 0 for immediate message sending
		Async:        true, // make the writer asynchronous
		RequiredAcks: kafka.RequireAll,
	}

	defer writer.Close()

	message := kafka.Message{
		Key:       []byte("123"),
		Value:     []byte("Hello, Kafka in Go!"),
		Time:      time.Now(),
		Partition: 2,
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Message sent successfully!")
	}
}

func consume() {
	brokers := []string{"localhost:9092"} // Replace with your broker addresses
	topic := "my-topic"                   // Replace with your desired topic
	groupID := "my-consumer-group"        // Replace with your consumer group ID

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	defer reader.Close()

	message, err := reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Received message:", string(message.Value))
	}
}
