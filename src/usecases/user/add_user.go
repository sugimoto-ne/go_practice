package usecases

import "github.com/sugimoto-ne/go_practice.git/domain"

type AddUser struct {
	UserRepository UserRepository
	// LastID         int64
}

func (au *AddUser) AddUser(name, password string) (int, error) {
	u := &domain.User{
		ID:       1,
		Name:     name,
		Password: password,
	}

	id, err := au.UserRepository.Store(u)
	if err != nil {
		return 0, err
	}

	return id, nil
}
