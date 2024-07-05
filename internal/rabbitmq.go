package internal

import (
	"fmt"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
 * RabbitClient is used to interact with RabbitMQ.
 * It reuses the TCP connection to the RabbitMQ server across the whole application.
 * For every concurrent task, it spawns new channels.
 * A channel is a multiplexed connection over the TCP connection.
 */
 type RabbitClient struct {
	conn *amqp.Connection // The TCP connection to the RabbitMQ server
	/*
	 * A multiplexed sub-connection to the RabbitMQ server.
	 * A channel is used to process/send messages.
	 */
	ch *amqp.Channel
}

// ConnectRabbitMQ establishes a connection to the RabbitMQ server using the provided credentials and vhost.
// username: the username for RabbitMQ authentication.
// password: the password for RabbitMQ authentication.
// host: the RabbitMQ server host and port.
// vhost: the virtual host to connect to on the RabbitMQ server.
func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

// NewRabbitMQClient creates a new RabbitClient with a channel from the provided connection.
// conn: the TCP connection to the RabbitMQ server.
func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	// Take the connection and spawn a new channel. This allows connection reuse.
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

// Close closes the RabbitMQ client's channel but not the connection.
func (rc RabbitClient) Close() error {
	return rc.ch.Close() // Close the channel, not the connection.
}

// CreateQueue declares a queue on the RabbitMQ server with the given parameters.
// queueName: the name of the queue.
// durable: if true, the queue will survive server restarts.
// autoDelete: if true, the queue will be deleted when no longer in use.
func (rc RabbitClient) CreateQueue(queueName string, durable bool, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(
		queueName,  // The name of the queue
		durable,    // If true, the queue will survive server restarts
		autoDelete, // If true, the queue will be deleted when no longer in use
		false,      // Exclusive - if true, the queue can only be used by the connection that declared it
		false,      // No-wait - if true, the server will not wait for the queue to be declared before responding
		nil,        // Arguments - additional parameters for queue declaration
	)
	return err
}

// CreateBinding creates a binding between a queue and an exchange with the given parameters.
// name: the name of the queue to bind to the exchange.
// binding: the routing key or pattern used for the binding.
// exchange: the name of the exchange to bind the queue to.
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving noWait=false makes the channel return an error if it fails to bind
	return rc.ch.QueueBind(
		name,     // The name of the queue to bind
		binding,  // The routing key or pattern for the binding
		exchange, // The name of the exchange to bind to
		false,    // noWait - if true, the server will not wait for the binding to be completed before responding
		nil,      // Arguments - additional parameters for binding declaration
	)
}


// Send publishes a message to a specified exchange with a given routing key.
//
// Parameters:
// - ctx: A context.Context for cancellation and timeouts.
// - exchange: The name of the exchange to publish the message to.
// - routingKey: The routing key to use for the message.
// - options: An amqp.Publishing struct containing the message properties and payload.
//
// The method uses PublishWithContext, which allows for context-based cancellation and timeouts.
// It sets both 'mandatory' and 'immediate' flags to false in the publish call.
//
// Returns an error if the publish operation fails, nil otherwise.
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(
		ctx,        // Context for cancellation and timeouts
		exchange,   // Exchange name
		routingKey, // Routing key
		true,      // Mandatory flag (if true, server returns an error if it cannot route the message)
		false,      // Immediate flag (if true, server returns an error if it cannot deliver to a consumer immediately)
		options,    // Publishing options including message properties and payload
	)
}