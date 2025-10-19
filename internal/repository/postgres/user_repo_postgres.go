package postgres

import (
	"context"
	"database/sql"

	"practice4-sqlx/internal/domain"
	"practice4-sqlx/internal/repository"

	"github.com/jmoiron/sqlx"
)

type userRepo struct{}

func NewUserRepo() repository.UserRepository { return &userRepo{} }

func (r *userRepo) Insert(ctx context.Context, ext sqlx.ExtContext, u domain.User) error {
	const q = `
INSERT INTO users (name, email, balance)
VALUES (:name, :email, :balance)
ON CONFLICT (email) DO NOTHING;
`
	_, err := sqlx.NamedExecContext(ctx, ext, q, u)
	return err
}

func (r *userRepo) GetAll(ctx context.Context, ext sqlx.ExtContext) ([]domain.User, error) {
	const q = `SELECT id, name, email, balance FROM users ORDER BY id`
	var out []domain.User
	if err := sqlx.SelectContext(ctx, ext, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *userRepo) GetByID(ctx context.Context, ext sqlx.ExtContext, id int) (domain.User, error) {
	const q = `SELECT id, name, email, balance FROM users WHERE id = $1`
	var u domain.User
	if err := sqlx.GetContext(ctx, ext, &u, q, id); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *userRepo) GetByIDForUpdate(ctx context.Context, ext sqlx.ExtContext, id int) (domain.User, error) {
	const q = `SELECT id, name, email, balance FROM users WHERE id = $1 FOR UPDATE`
	var u domain.User
	if err := sqlx.GetContext(ctx, ext, &u, q, id); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *userRepo) AddBalance(ctx context.Context, ext sqlx.ExtContext, id int, delta float64) error {
	const q = `UPDATE users SET balance = balance + $1 WHERE id = $2`
	res, err := ext.ExecContext(ctx, q, delta, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

var _ repository.UserRepository = (*userRepo)(nil)
