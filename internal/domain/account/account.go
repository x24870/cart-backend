package account

import (
	"cart-backend/internal/domain"
	"context"
)

type Account struct {
	domain.Base
	Address string `gorm:"column:address;type:varchar(42);primary_key" json:"address"`
}

type Repository interface {
	FirstOrCreate(ctx context.Context, address string) (*Account, error)
	GetByAddress(ctx context.Context, address string) (*Account, error)
}

// table name
func (Account) TableName() string {
	return "account"
}
