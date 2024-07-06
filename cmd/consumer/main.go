package main

import (
	"log" // Package for logging

	"github.com/mapfumo/event-driven-rabbitmq/internal" // Internal package for RabbitMQ utilities
)

func main() {
	// Establish a connection to the RabbitMQ server with the provided credentials and vhost
	conn, err := internal.ConnectRabbitMQ("tony", "secret", "localhost:5672", "customers")
	if err != nil {
		// If there's an error while connecting, panic
		panic(err)
	}

	// Create a new RabbitMQ client using the established connection
	mqClient, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		// If there's an error while creating the client, panic
		panic(err)
	}

	// Start consuming messages from the "customers_created" queue with the consumer name "email-service"
	// autoAck is set to false, meaning messages need to be acknowledged manually
	messageBus, err := mqClient.Consume("customers_created", "email-service", false)
	if err != nil {
		// If there's an error while setting up the consumer, panic
		panic(err)
	}

	// Create a blocking channel to keep the main function running indefinitely
	var blocking chan struct{}

	// Start a goroutine to process messages from the messageBus channel
	go func() {
		for message := range messageBus {
			// Log the received message body
			log.Printf("New Message: %v", string(message.Body))
		}
	}()

	// Log that the consumer is running and how to close the program
	log.Println("Consuming, to close the program press CTRL+C")

	// Block forever, keeping the main function running
	<-blocking
}
