package main

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/jsonschema"
)

// Purchase is a simple record example
type Purchase struct {
	Id           int    `json:"id"`
	Quantity     int    `json:"quantity"`
	ItemType     string `json:"item_type"`
	PricePerUnit string `json:"price_per_unit"`
}

func main() {
	// Get env
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bootstrap-servers> <group> <topics..>\n", os.Args[0])
		os.Exit(1)
	}

	bootstrapServers := os.Args[1]
	clusterApiKey := os.Args[2]
	clusterApiSecret := os.Args[3]
	topic := os.Args[4]
	schemaRegistryUrl := os.Args[5]
	schemaRegistryApiKey := os.Args[6]
	schemaRegistryApiSecret := os.Args[7]

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to get hostname: %s", err)
		os.Exit(1)
	}

	// Create producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"client.id":         hostname,
		"sasl.username":     clusterApiKey,
		"sasl.password":     clusterApiSecret,
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	// Create schema register client
	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		schemaRegistryUrl,
		schemaRegistryApiKey,
		schemaRegistryApiSecret,
	))
	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	ser, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	// Serialize message
	msg := Purchase{
		Id:           5950,
		Quantity:     3,
		ItemType:     "guitar",
		PricePerUnit: "BMg=",
	}
	msgBytes, err := ser.Serialize(topic, &msg)
	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	// Asynchronous writes
	deliveryChan := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msgBytes,
		Key:            []byte(time.Now().UTC().Format(time.RFC1123)),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Failed to deliver message: %s", err)
		os.Exit(1)
	}

	// Go-routine to handle message delivery reports and possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	// Wait for all messages to be delivered
	// Flush and close the producer and the events channel
	for p.Flush(1000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n")
	}
	p.Close()
}
