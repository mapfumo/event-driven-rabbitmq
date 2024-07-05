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

## Tests

```bash
go test -v ./...
```

## References

- [RabbitMQ](https://www.rabbitmq.com/)
- [RabbitMQ GoLang Client](https://github.com/rabbitmq/amqp091-go)
- [RabbitMQ GoLang Client Examples](https://github.com/rabbitmq/rabbitmq-tutorials/tree/main/go)
