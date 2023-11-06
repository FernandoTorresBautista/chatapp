package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Rabbit structure
type Rabbit struct {
	user     string
	password string
	logger   *log.Logger
	Rmq      *amqp.Connection
	rooms    []string
}

// NewRabbit ...
func NewRabbit(logger *log.Logger, user, password string) *Rabbit {
	return &Rabbit{
		user:     user,
		password: password,
		logger:   logger,
		Rmq:      nil,
		rooms:    []string{},
	}
}

// Start rabbit ...
func (mq *Rabbit) Start() error {
	var err error
	mq.Rmq, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", mq.user, mq.password))
	if err != nil {
		mq.logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return err
	}
	return nil
}

// Stop rabbit ...
func (mq *Rabbit) Stop() error {
	return mq.Rmq.Close()
}

// CreateQueue ...
func (mq *Rabbit) CreateQueue(room string) error {
	ch, err := mq.Rmq.Channel()
	if err != nil {
		mq.logger.Fatalf("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	mq.rooms = append(mq.rooms, room)

	queueName := room
	arguments := amqp.Table{
		"x-max-length": 50, // 50 messages limit
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		arguments, // Arguments
	)
	if err != nil {
		mq.logger.Fatalf("Failed to declare a queue: %v", err)
		return err
	}
	return nil
}

// SendMessage ...
func (mq *Rabbit) SendMessage(mt int, room, msg string) error {
	ch, err := mq.Rmq.Channel()
	if err != nil {
		mq.logger.Fatalf("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	queueName := room
	err = ch.Publish(
		"",
		queueName,
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf("Msg sended: %s to the room: %s", msg, room)
	return nil
}
