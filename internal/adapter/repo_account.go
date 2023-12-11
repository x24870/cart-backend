package adapter

import (
	"cart-backend/internal/domain/account"
	"context"

	"gorm.io/gorm"
)

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) account.Repository {
	return &accountRepo{db: db}
}

func (r *accountRepo) CreateIfNotExists(ctx context.Context, account *account.Account) error {
	return nil
}

func (r *accountRepo) GetByAddress(ctx context.Context, address string) (*account.Account, error) {
	return nil, nil
}
