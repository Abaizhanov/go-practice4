package repository

import (
	"context"

	"practice4-sqlx/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Insert(ctx context.Context, ext sqlx.ExtContext, u domain.User) error
	GetAll(ctx context.Context, ext sqlx.ExtContext) ([]domain.User, error)
	GetByID(ctx context.Context, ext sqlx.ExtContext, id int) (domain.User, error)
	GetByIDForUpdate(ctx context.Context, ext sqlx.ExtContext, id int) (domain.User, error)
	AddBalance(ctx context.Context, ext sqlx.ExtContext, id int, delta float64) error
}
