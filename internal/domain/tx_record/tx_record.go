package txrecord

import (
	"cart-backend/internal/domain"
	"context"
)

// TxRecord ...
type TxRecord struct {
	domain.Base
	Address     *string `gorm:"column:address;type:varchar(42)" json:"address"`
	ProjectName string  `gorm:"column:project_name;type:varchar(255)" json:"project_name"`
	Url         string  `gorm:"column:url;type:varchar(2048)" json:"url"`
	Amount      string  `gorm:"column:amount;type:varchar(255)" json:"amount"`
	Symbol      string  `gorm:"column:symbol;type:varchar(255)" json:"symbol"`
}

type Repository interface {
	Create(ctx context.Context, txRecord *TxRecord) error
	ListByAddress(ctx context.Context, address string) (*[]TxRecord, error)
}
