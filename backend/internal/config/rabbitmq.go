package config

import (
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	URL string
}

func LoadRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		URL: getRabbitMQEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func ConnectRabbitMQ(config *RabbitMQConfig) (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SetupStockQueue(ch *amqp.Channel) error {
	// Declare the stock queue
	_, err := ch.QueueDeclare(
		"stock_commands", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	// Declare the stock responses queue
	_, err = ch.QueueDeclare(
		"stock_responses", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func getRabbitMQEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
