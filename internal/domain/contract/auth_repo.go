package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByID(ctx context.Context, id int64) (*entity.User, error)
}
