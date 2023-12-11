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
}

type Repository interface {
	Create(ctx context.Context, txRecord *TxRecord) error
	ListByAddress(ctx context.Context, address string) (*[]TxRecord, error)
}
