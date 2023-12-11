package adapter

import (
	txrecord "cart-backend/internal/domain/tx_record"
	"context"

	"gorm.io/gorm"
)

type txRecordRepo struct {
	db *gorm.DB
}

func NewTxRecordRepo(db *gorm.DB) txrecord.Repository {
	return &txRecordRepo{db: db}
}

func (r *txRecordRepo) Create(ctx context.Context, txRecord *txrecord.TxRecord) error {
	return r.db.WithContext(ctx).Create(txRecord).Error
}

func (r *txRecordRepo) ListByAddress(ctx context.Context, address string) (*[]txrecord.TxRecord, error) {
	var txRecords []txrecord.TxRecord
	err := r.db.WithContext(ctx).Where("address = ?", address).Find(&txRecords).Error
	if err != nil {
		return nil, err
	}
	return &txRecords, nil
}
