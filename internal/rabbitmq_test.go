package internal

import (
	"testing"
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