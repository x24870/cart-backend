package txrecord

import (
	"cart-backend/internal/domain"
	"cart-backend/internal/domain/account"
	"context"
)

// TxRecord ...
type TxRecord struct {
	domain.Base
	Account     account.Account `gorm:"foreignKey:Address"`
	ProjectName string          `gorm:"column:project_name;type:varchar(255)" json:"project_name"`
	Url         string          `gorm:"column:url;type:varchar(2048)" json:"url"`
	Amount      string          `gorm:"column:amount;type:varchar(255)" json:"amount"`
	Symbol      string          `gorm:"column:symbol;type:varchar(255)" json:"symbol"`
}

type Repository interface {
	Create(ctx context.Context, txRecord *TxRecord) error
	ListByAddress(ctx context.Context, address string) (*[]TxRecord, error)
}
