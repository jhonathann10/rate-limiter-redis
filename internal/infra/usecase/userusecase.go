package usecase

import (
	"context"

	"github.com/jhonathann10/rate-limiter-redis/internal/entity"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/internalerrors"
)

type UserOutputDTO struct {
	Username string `json:"username"`
}

type UserUseCase struct {
	userRepository entity.UserRepositoryInterface
}

type UserUseCaseInterface interface {
	SaveUser(ctx context.Context, user string) *internalerrors.InternalError
	GetUser(ctx context.Context) (*UserOutputDTO, *internalerrors.InternalError)
}

func NewUserUseCase(userRepository entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

func (u *UserUseCase) SaveUser(ctx context.Context, user string) *internalerrors.InternalError {
	errSaveUser := u.userRepository.SaveUser(ctx, user)
	if errSaveUser != nil {
		return errSaveUser
	}

	return nil
}

func (u *UserUseCase) GetUser(ctx context.Context) (*UserOutputDTO, *internalerrors.InternalError) {
	user, errGetUser := u.userRepository.GetUser(ctx)
	if errGetUser != nil {
		return nil, errGetUser
	}

	return &UserOutputDTO{Username: user.Username}, nil
}
