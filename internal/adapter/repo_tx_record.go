package adapter

import (
	t "cart-backend/internal/domain/transaction"
	utils "cart-backend/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type txRecordRepo struct {
	db *gorm.DB
}

func NewTxRecordRepo(db *gorm.DB) t.TxRecordRepo {
	return &txRecordRepo{db: db}
}

func (r *txRecordRepo) Create(
	ctx context.Context, tx *t.TxRecord, ops *[]t.Operation, its *[]t.Intent,
) error {
	if err := utils.Transactional(r.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(tx).Error; err != nil {
			return err
		}
		for _, op := range *ops {
			if err := tx.WithContext(ctx).Create(&op).Error; err != nil {
				return err
			}
		}
		for _, it := range *its {
			if err := tx.WithContext(ctx).Create(&it).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *txRecordRepo) ListByAddress(
	ctx context.Context, address string,
) (*[]t.TxRecord, error) {
	var txRecords []t.TxRecord
	err := r.db.WithContext(ctx).Where("address = ?", address).Find(&txRecords).Error
	if err != nil {
		return nil, err
	}
	return &txRecords, nil
}
