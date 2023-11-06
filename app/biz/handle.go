package biz

import (
	"log"

	"chatapp/app/biz/models"
	"chatapp/app/client/db"
	"chatapp/app/client/rabbitmq"
)

// Biz structure with all escential elements
type Biz struct {
	logger *log.Logger
	db     db.Repository
	Rabbit *rabbitmq.Rabbit
}

// New return a new biz instance
func New(logger *log.Logger, db db.Repository) *Biz {
	return &Biz{
		logger: logger,
		db:     db,
	}
}

// SetRabbit set the rabbitmq instance
func (b *Biz) SetRabbit(rmq *rabbitmq.Rabbit) {
	b.Rabbit = rmq
}

// Handle implementations ...
type Handle interface {
	// db
	CreateUser(nu *models.User) (uint64, error)
	GetUser(email, pass string) (*models.User, error)

	// rabbitmq
	CreateQueue(room string) error
	SendMessage(mt int, room, msg string) error
}

// Start all elements of the biz
func (b *Biz) Start() error {
	// initialize the Biz
	if b.db != nil {
		err := b.db.Start()
		if err != nil {
			return err
		}
	}
	if b.Rabbit != nil {
		err := b.Rabbit.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop all elements of the biz
func (b *Biz) Stop() error {
	// Stop the Biz
	if b.db != nil {
		err := b.db.Stop()
		if err != nil {
			return err
		}
	}
	if b.Rabbit != nil {
		err := b.Rabbit.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// implement functions

// CreateUser with the repository, we can check here any other concern with the ew user...
func (b *Biz) CreateUser(nu *models.User) (uint64, error) {
	return b.db.CreateUser(nu)
}

// GetUser search a user in db
func (b *Biz) GetUser(email, pass string) (*models.User, error) {
	return b.db.GetUser(email, pass)
}

// CreateQueue ...
func (b *Biz) CreateQueue(room string) error {
	return b.Rabbit.CreateQueue(room)
}

// SendMessage ...
func (b *Biz) SendMessage(mt int, room, msg string) error {
	return b.Rabbit.SendMessage(mt, room, msg)
}
