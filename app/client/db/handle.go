package db

import "chatapp/app/biz/models"

type Repository interface {
	// methods of the clients/[db, kafka, redis, etc..]
	Start() error
	Stop() error

	GetExample() (interface{}, error)
	CreateExample(fiel string) (interface{}, error)

	CreateUser(nu *models.User) (uint64, error)
	GetUser(email, pass string) (*models.User, error)
}
