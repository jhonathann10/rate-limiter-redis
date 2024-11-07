package entity

import (
	"context"

	"github.com/jhonathann10/rate-limiter-redis/internal/infra/internalerrors"
)

type User struct {
	Username string `json:"username"`
}

type UserRepositoryInterface interface {
	SaveUser(ctx context.Context, user string) *internalerrors.InternalError
	GetUser(ctx context.Context) (*User, *internalerrors.InternalError)
}
