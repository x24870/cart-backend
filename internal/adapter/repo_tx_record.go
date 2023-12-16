package adapter

import (
	t "cart-backend/internal/domain/transaction"
	"context"

	"gorm.io/gorm"
)

type txRecordRepo struct {
	db *gorm.DB
}

func NewTxRecordRepo(db *gorm.DB) t.TxRecordRepo {
	return &txRecordRepo{db: db}
}

func (r *txRecordRepo) Create(ctx context.Context, txRecord *t.TxRecord) error {
	return r.db.WithContext(ctx).Create(txRecord).Error
}

func (r *txRecordRepo) ListByAddress(ctx context.Context, address string) (*[]t.TxRecord, error) {
	var txRecords []t.TxRecord
	err := r.db.WithContext(ctx).Where("address = ?", address).Find(&txRecords).Error
	if err != nil {
		return nil, err
	}
	return &txRecords, nil
}
