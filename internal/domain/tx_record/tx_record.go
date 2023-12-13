package txrecord

import (
	"cart-backend/internal/domain"
	"context"
)

// TxRecord ...
type TxRecord struct {
	domain.Base
	Address     string `gorm:"column:address;type:varchar(42);index;reference:Address"` // Adjusted for FK relationship
	ProjectName string `gorm:"column:project_name;type:varchar(255)"`
	Url         string `gorm:"column:url;type:varchar(2048)"`
	Amount      string `gorm:"column:amount;type:varchar(255)"`
	Symbol      string `gorm:"column:symbol;type:varchar(255)"`
	Signature   string `gorm:"column:signature;type:varchar(255)"`
}

type Repository interface {
	Create(ctx context.Context, txRecord *TxRecord) error
	ListByAddress(ctx context.Context, address string) (*[]TxRecord, error)
}

// table name
func (TxRecord) TableName() string {
	return "tx_record"
}
