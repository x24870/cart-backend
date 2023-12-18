package adapter

import (
	t "cart-backend/internal/domain/transaction"
	utils "cart-backend/pkg/utils"
	"context"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type txRecordRepo struct {
	db *gorm.DB
}

func NewTxRecordRepo(db *gorm.DB) t.TxRecordRepo {
	return &txRecordRepo{db: db}
}

func (r *txRecordRepo) Create(
	ctx context.Context, t *t.TxRecord, ops *[]t.Operation, its *[]t.Intent,
) error {
	if err := utils.Transactional(r.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&t).Error; err != nil {
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
	err := r.db.WithContext(ctx).Where("account = ?", address).Order("created_at desc").Find(&txRecords).Error
	if err != nil {
		return nil, err
	}
	return &txRecords, nil
}

type operationRepo struct {
	db *gorm.DB
}

func NewOperationRepo(db *gorm.DB) t.OperationRepo {
	return &operationRepo{db: db}
}

func (o *operationRepo) ListByTxHash(
	ctx context.Context, txHash string,
) (*[]t.Operation, error) {
	var operations []t.Operation
	err := o.db.WithContext(ctx).Where("tx_hash = ?", txHash).Order("created_at").Find(&operations).Error
	if err != nil {
		return nil, err
	}
	return &operations, nil
}

type intentRepo struct {
	db *gorm.DB
}

func NewIntentRepo(db *gorm.DB) t.IntentRepo {
	return &intentRepo{db: db}
}

func (i *intentRepo) ListByOperationID(
	ctx context.Context, operationID uuid.UUID,
) (*[]t.Intent, error) {
	var intents []t.Intent
	err := i.db.WithContext(ctx).Where("operation_id = ?", operationID).Order("created_at").Find(&intents).Error
	if err != nil {
		return nil, err
	}
	return &intents, nil
}
