package internal

import (
	"testing"
	amqp "github.com/rabbitmq/amqp091-go"
	"context"
	"time"
)


func TestConnectRabbitMQ(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"

	conn, err := ConnectRabbitMQ(username, password, host, vhost)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("Expected a valid connection, got nil")
	}
}

func TestNewRabbitMQClient(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"

	conn, err := ConnectRabbitMQ(username, password, host, vhost)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	rc, err := NewRabbitMQClient(conn)
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rc.Close()

	if rc.ch == nil {
		t.Fatal("Expected a valid channel, got nil")
	}
}

func TestCreateQueue(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"

	conn, err := ConnectRabbitMQ(username, password, host, vhost)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	rc, err := NewRabbitMQClient(conn)
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rc.Close()

	// create durable queue
	err = rc.CreateQueue("customers_created", true, false)
	if err != nil {
		t.Fatalf("Failed to create queue: %v", err)
	}

		// create non-durable queue. Should disappear after restarting rabbitmq
		// docker restart rabbitmq
		err = rc.CreateQueue("customers_not_durable", false, false)
		if err != nil {
			t.Fatalf("Failed to create queue: %v", err)
		}
}

func TestCreateBinding(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"
	queueName := "customers_created"
	binding := "customers.created.*"
	exchange := "customers_events"

	// Establish connection to RabbitMQ
	conn, err := ConnectRabbitMQ(username, password, host, vhost)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a RabbitMQ client
	rc, err := NewRabbitMQClient(conn)
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rc.Close()

	// Declare the queue
	err = rc.CreateQueue(queueName, true, false)
	if err != nil {
		t.Fatalf("Failed to create queue: %v", err)
	}

	// Create the binding between the queue and the pre-existing exchange
	err = rc.CreateBinding(queueName, binding, exchange)
	if err != nil {
		t.Fatalf("Failed to create binding: %v", err)
	}
}

func TestSend(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"
	exchange := "customers_events"
	routingKey := "customers.created.us"
	options := amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent, //also there's amqp.Transient
		Body: []byte("A cool message between services"),
	}

	// Establish connection to RabbitMQ
	conn, err := ConnectRabbitMQ(username, password, host, vhost)
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a RabbitMQ client
	rc, err := NewRabbitMQClient(conn)
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ client: %v", err)
	}
	defer rc.Close()

	// Create a background context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send a message
	err = rc.Send(ctx, exchange, routingKey, options)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
}