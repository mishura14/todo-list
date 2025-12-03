package repository

import "git-register-project/internal/models"

type UserRepository interface {
	CheckEmailExists(email string) (bool, error)
	CreateUser(user *models.UserRedis) error
}
