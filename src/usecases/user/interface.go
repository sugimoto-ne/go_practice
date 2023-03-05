package usecases

import "github.com/sugimoto-ne/go_practice.git/domain"

type UserRepository interface {
	FindById(id domain.UserID) (*domain.User, error)
	FindAll() (domain.Users, error)
	Store(u *domain.User) (int, error)
}
