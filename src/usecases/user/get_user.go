package usecases

import "github.com/sugimoto-ne/go_practice.git/domain"

type GetUser struct {
	UserRepository UserRepository
}

func (gu *GetUser) GetUser(id domain.UserID) (*domain.User, error) {
	u, err := gu.UserRepository.FindById(id)

	if err != nil {
		return nil, err
	}

	return u, nil
}
