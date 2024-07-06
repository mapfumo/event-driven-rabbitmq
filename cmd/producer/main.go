package main

import (
	"context"
	"log"
	"time"

	"github.com/mapfumo/event-driven-rabbitmq/internal"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// Establish a connection to the RabbitMQ server
	conn, err := internal.ConnectRabbitMQ("tony", "secret", "localhost:5672", "customers")
	if err != nil {
		// If there's an error while connecting, panic
		panic(err)
	}
	defer conn.Close() // Ensure the connection is closed when the main function exits

	// Create a new RabbitMQ client using the established connection
	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		// If there's an error while creating the client, panic
		panic(err)
	}
	defer client.Close() // Ensure the client's channel is closed when the main function exits

	// Create a durable queue named "customers_created"
	if err := client.CreateQueue("customers_created", true, false); err != nil {
		// If there's an error while creating the queue, panic
		panic(err)
	}

	// Create a non-durable, auto-delete queue named "customers_test"
	if err := client.CreateQueue("customers_test", false, true); err != nil {
		// If there's an error while creating the queue, panic
		panic(err)
	}

	log.Println(client) // Log the client details

	// Create a binding between the "customers_events" exchange and the "customers_created" queue
	if err := client.CreateBinding("customers_created", "customers.created.*", "customers_events"); err != nil {
		// If there's an error while creating the binding, panic
		panic(err)
	}

	// Create a binding between the "customers_events" exchange and the "customers_test" queue
	if err := client.CreateBinding("customers_test", "customers.*", "customers_events"); err != nil {
		// If there's an error while creating the binding, panic
		panic(err)
	}

	// Create a context with a 5-second timeout to manage message publishing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the context is canceled when the main function exits

	// Send a persistent message to the "customers_events" exchange with the routing key "customers.created.se"
	if err := client.Send(ctx, "customers_events", "customers.created.se", amqp091.Publishing{
		ContentType:  "text/plain",       // The message payload is plain text
		DeliveryMode: amqp091.Persistent, // The message should be saved (durable) if no resources accept it before a restart
		Body:         []byte("A cool message between services"),
	}); err != nil {
		// If there's an error while sending the message, panic
		panic(err)
	}

	// Send a transient message to the "customers_events" exchange with the routing key "customers.test"
	if err := client.Send(ctx, "customers_events", "customers.test", amqp091.Publishing{
		ContentType:  "text/plain",      // The message payload is plain text
		DeliveryMode: amqp091.Transient, // The message can be deleted (non-durable) if no resources accept it before a restart
		Body:         []byte("A second cool message"),
	}); err != nil {
		// If there's an error while sending the message, panic
		panic(err)
	}

	log.Println(client) // Log the client details again
}
