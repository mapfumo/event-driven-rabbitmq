package internal

import (
	"context"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

func TestConsume(t *testing.T) {
	username := "tony"
	password := "secret"
	host := "localhost:5672"
	vhost := "customers"
	queue := "customers_created"
	consumer := "email_service"
	autoAck := false

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

	// Declare the queue (assuming it may not exist)
	err = rc.CreateQueue(queue, true, false)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// Publish a test message to the queue
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testMessage := amqp.Publishing{
		ContentType: "text/plain",
		DeliveryMode: amqp.Persistent,
		Body: []byte("A cool message between services"),
	}

	err = rc.Send(ctx, "", queue, testMessage)
	if err != nil {
		t.Fatalf("Failed to send test message: %v", err)
	}

	// Consume messages from the queue
	messages, err := rc.Consume(queue, consumer, autoAck)
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	// Read the first message from the channel
	select {
	case msg := <-messages:
		if string(msg.Body) != "A cool message between services" {
			t.Fatalf("Expected message body 'Test message', got '%s'", string(msg.Body))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for message")
	}
}