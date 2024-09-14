package services

import (
	"github.com/Ferdinand-work/PetalPix/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	CreateUser(*models.User) (*mongo.InsertOneResult, error)
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
	Follow(interface{}, string) (*[]string, error)
	GetFollowing(string) (*[]models.User, error)
	Unfollow(interface{}, string) (*[]string, error)
}
