package rabbitmq

import (
	"concurrency/src/utils"
	"os"

	"github.com/streadway/amqp"
)

func ConnectToRabbit() (*amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError(err, "failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to open a channel")

	utils.LogWithInfo("connected to rabbitMQ", "rabbitmq")
	return ch, err
}

type RabbitMQClient struct {
	ch   *amqp.Channel
	name string
}

func CreateRabbitMQClient(r *amqp.Channel, name string) *RabbitMQClient {
	return &RabbitMQClient{
		ch:   r,
		name: name,
	}
}

func (rmq *RabbitMQClient) CreateRabbitMQueue() (amqp.Queue, error) {
	q, err := rmq.ch.QueueDeclare(
		rmq.name, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	utils.FailOnError(err, "failed to declare a queue")
	return q, err
}

func (rmq *RabbitMQClient) PublishMessage(q amqp.Queue, body []byte) error {

	err := rmq.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	utils.FailOnError(err, "Failed to publish a message")
	return err
}
