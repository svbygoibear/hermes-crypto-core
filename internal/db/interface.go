package db

import "hermes-crypto-core/internal/models"

type DBInterface interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	UpdateUser(id string, user models.User, updateScore bool) (*models.User, error)
	DeleteUser(id string) error
}

var DB DBInterface
