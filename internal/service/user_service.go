package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"practice4-sqlx/internal/domain"
	"practice4-sqlx/internal/repository"
)

type UserService struct {
	DB   *sqlx.DB
	Repo repository.UserRepository
}

func NewUserService(db *sqlx.DB, repo repository.UserRepository) *UserService {
	return &UserService{DB: db, Repo: repo}
}

func (s *UserService) InsertUser(ctx context.Context, name, email string, balance float64) error {
	u := domain.User{
		Name:    name,
		Email:   email,
		Balance: balance,
	}
	return s.Repo.Insert(ctx, s.DB, u)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.Repo.GetAll(ctx, s.DB)
}

func (s *UserService) TransferBalance(ctx context.Context, fromID, toID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if fromID == toID {
		return errors.New("cannot transfer to the same user")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	from, err := s.Repo.GetByIDForUpdate(ctx, tx, fromID)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}
	if from.Balance < amount {
		return fmt.Errorf("insufficient funds: have %.2f need %.2f", from.Balance, amount)
	}
	if _, err := s.Repo.GetByIDForUpdate(ctx, tx, toID); err != nil {
		return fmt.Errorf("receiver not found: %w", err)
	}

	if err := s.Repo.AddBalance(ctx, tx, fromID, -amount); err != nil {
		return fmt.Errorf("debit failed: %w", err)
	}
	if err := s.Repo.AddBalance(ctx, tx, toID, amount); err != nil {
		return fmt.Errorf("credit failed: %w", err)
	}

	return tx.Commit()
}
