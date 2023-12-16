package transaction

import (
	"cart-backend/internal/domain"
	"cart-backend/internal/domain/common"
	"context"

	uuid "github.com/satori/go.uuid"
)

// TxRecord ...
type TxRecord struct {
	domain.Base
	Hash      string `gorm:"column:hash;type:varchar(255);primary_key"`
	Account   string `gorm:"column:account;type:varchar(42);foreignKey;index;reference:Address"` // Adjusted for FK relationship
	Signature string `gorm:"column:signature;type:varchar(255)"`
}

func (TxRecord) TableName() string {
	return "tx_record"
}

// Operation ...
type Operation struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
	Tx          TxRecord  `gorm:"foreignKey:Hash;references:Hash"` // Foreign key relationship
	ProjectName string    `gorm:"column:project_name;type:varchar(255)"`
	Url         string    `gorm:"column:url;type:varchar(2048)"`
	common.Token
}

func (Operation) TableName() string {
	return "operation"
}

// Intent ...
type Intent struct {
	Operation   Operation `gorm:"foreignKey:ID;references:ID"` // Foreign key relationship
	Description string    `gorm:"column:description;type:varchar(255)"`
}

func (Intent) TableName() string {
	return "intent"
}

type TxRecordRepo interface {
	Create(ctx context.Context, txRecord *TxRecord, ops *[]Operation, intents *[]Intent) error
	ListByAddress(ctx context.Context, address string) (*[]TxRecord, error)
}
