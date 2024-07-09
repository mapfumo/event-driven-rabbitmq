# Event-Driven Architecture (EDA) - With RabbitMQ & GoLang

## Introduction

This is a simple example of how to implement an event-driven architecture (EDA) using RabbitMQ as the message broker.

## Prerequisites

- Docker
- GoLang

## Docker and RabbitMQ Setup

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management

```

- go to http://localhost:15672 and login with `guest` and `guest` credentials
- Add your own credentials and remove default guest login

```bash
alias rabbitmqctl='docker exec -it rabbitmq rabbitmqctl'
rabbitmqctl add_user tony secret
rabbitmqctl set_user_tags tony administrator
rabbitmqctl delete_user guest
docker logs rabbitmq --follow
```

## Virtual Hosts

RabbitMQ virtual hosts (vhosts) provide logical separation within a single RabbitMQ instance, allowing different applications to have isolated queues, exchanges, and bindings. They also enable fine-grained access control by granting specific permissions to users for different vhosts.

```bash
rabbitmqctl add_vhost customers`` ## going to hold all our future resources related to customers
```

Give config, write and read permissions to user `tony` for vhost `customers`

```bash
rabbitmqctl set_permissions -p customers tony ".*" ".*" ".*"
```

## Exchanges

- Exchanges are like mailboxes. They are used to route messages to queues. Exchanges can be configured to be durable or transient.
- Set permissions for user `tony` for the `customers_events` exchange

```bash
rabbitmqctl set_topic_permissions -p customers tony customer_events "^customers.*" "^customers.*"
```

````bash
alias rabbitmqadmin='docker exec -it rabbitmq rabbitmqadmin'
rabbitmqadmin declare exchange --vhost=customers name=customers_events type=topic -u tony -p secret durable=true
## Tests

```bash
go test -v ./...
````

## References

- [RabbitMQ](https://www.rabbitmq.com/)
- [RabbitMQ GoLang Client](https://github.com/rabbitmq/amqp091-go)
- [RabbitMQ GoLang Client Examples](https://github.com/rabbitmq/rabbitmq-tutorials/tree/main/go)
- [Learn RabbitMQ for Event-Driven Architecture](https://programmingpercy.tech/blog/event-driven-architecture-using-rabbitmq/)
