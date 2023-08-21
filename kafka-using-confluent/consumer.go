package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <bootstrap-servers> <group> <topics..>\n", os.Args[0])
		os.Exit(1)
	}

	bootstrapServers := os.Args[1]
	clusterApiKey := os.Args[2]
	clusterApiSecret := os.Args[3]
	group := os.Args[4]
	topics := os.Args[5:]

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"sasl.mechanisms":       "PLAIN",
		"security.protocol":     "SASL_SSL",
		"sasl.username":         clusterApiKey,
		"sasl.password":         clusterApiSecret,
		"group.id":              group,
		"session.timeout.ms":    6000,
		// Start reading from the first message of each assigned
		// partition if there are no previously committed offsets
		// for this group.
		"auto.offset.reset": "earliest",
		// Whether we store offsets automatically.
		"enable.auto.offset.store": false,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	if err := c.SubscribeTopics(topics, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe topics: %s\n", err)
		os.Exit(1)
	}

	run := true

	for run {
		select {
		case sig := <-sigChan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			// Retrieve records one-by-one
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				// We can store the offsets of the messages manually or let
				// the library do it automatically based on the setting
				// enable.auto.offset.store. Once an offset is stored, the
				// library takes care of periodically committing it to the broker
				// if enable.auto.commit isn't set to false (the default is true).
				// By storing the offsets manually after completely processing
				// each message, we can ensure atleast once processing.
				if _, err := c.StoreMessage(e); err != nil {
					fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n", e.TopicPartition)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to automatically recover.
				// But in this example we choose to terminate the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}
