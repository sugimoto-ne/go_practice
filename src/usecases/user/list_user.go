package usecases

import "github.com/sugimoto-ne/go_practice.git/domain"

type ListUser struct {
	UserRepository UserRepository
}

func (lu *ListUser) ListUser() (domain.Users, error) {
	users, err := lu.UserRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
