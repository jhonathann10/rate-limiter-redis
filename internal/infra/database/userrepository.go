package database

import (
	"context"
	"net/http"
	"time"

	"github.com/jhonathann10/rate-limiter-redis/internal/entity"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/internalerrors"
	"github.com/redis/go-redis/v9"
)

const keyUser = "user"

type UserRepository struct {
	client *redis.Client
}

func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (ur *UserRepository) SaveUser(ctx context.Context, user string) *internalerrors.InternalError {
	err := ur.client.Set(ctx, keyUser, user, 10*time.Second).Err()
	if err != nil {
		return &internalerrors.InternalError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}

func (ur *UserRepository) GetUser(ctx context.Context) (*entity.User, *internalerrors.InternalError) {
	username, err := ur.client.Get(ctx, keyUser).Result()
	if err != nil {
		return nil, &internalerrors.InternalError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &entity.User{Username: username}, nil
}
