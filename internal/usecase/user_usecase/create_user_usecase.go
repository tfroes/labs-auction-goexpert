package user_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
)

type CreateUserInputDTO struct {
	Name string `json:"name" binding:"required,min=1"`
}

func (u *UserUseCase) CreateUser(
	ctx context.Context,
	createUserInput CreateUserInputDTO) *internal_error.InternalError {
	user, err := user_entity.CreateUser(
		createUserInput.Name)

	if err != nil {
		return err
	}

	if err := u.UserRepository.CreateUser(
		ctx, user); err != nil {
		return err
	}

	return nil
}
