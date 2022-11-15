package rabbitmq

import (
	"concurrency/src/utils"
	"os"

	"github.com/streadway/amqp"
)

func ConnectToRabbit() (*amqp.Channel, error) {
	addr := os.Getenv("RABBITMQ_ADDRESS")
	conn, err := amqp.Dial(addr)
	utils.FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "failed to open a channel")
	defer ch.Close()

	// We create a Queue to send the message to.
	utils.LogWithInfo("connected to rabbitMQ", "rabbitmq")

	return ch, err
}

func CreateRabbitMQueue(r *amqp.Channel, name string) (amqp.Queue, error) {
	q, err := r.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "failed to declare a queue")
	return q, err
}

func PublishMessage(name string, ch *amqp.Channel, body []byte) {

	err := ch.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	utils.FailOnError(err, "Failed to publish a message")
}
