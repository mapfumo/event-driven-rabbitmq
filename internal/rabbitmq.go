package internal

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

/*
 * Reuse the TCP connection to RabbitMQ server accross the whole application
 * Spawn new channels on every concurrent task that is running
 * A channel is a multiplexed connectiion over the TCP connetion
 */
type RabbitClient struct {
	conn *amqp.Connection // The TCP connection to RabbitMQ server
	/*
	 * multioplexed sub-connection to RabbitMQ server
	 * A channel is used to process/send messages
	 */
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))

}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	// take the connection and spawn a new channel. This allows connection reuse
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close() // close the channel, not the connection
}

func (rc RabbitClient) CreateQueue(queueName string, durable bool, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(
		queueName,  // name
		durable,    // durable
		autoDelete, // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	return err
}
