package store

import (
	"errors"

	"github.com/sugimoto-ne/go_practice.git/domain"
)

type UserRepository struct {
	// Repository
	LastID domain.UserID
}

func (ur *UserRepository) Store(u *domain.User) (int, error) {
	ur.LastID++
	u.ID = domain.UserID(ur.LastID)
	Users = append(Users, *u)

	return len(Users), nil
}

func (ur *UserRepository) FindById(id domain.UserID) (*domain.User, error) {

	if len(Users) >= int(id) {
		return &Users[id-1], nil
	}

	return nil, errors.New("not found")
}

func (ur *UserRepository) FindAll() (domain.Users, error) {
	return Users, nil
}
