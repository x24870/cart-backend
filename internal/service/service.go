package service

import (
	"cart-backend/internal/domain/account"
	t "cart-backend/internal/domain/transaction"
	"context"
	"fmt"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*CreateResponse, error)
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
}

type service struct {
	accountRepo  account.Repository
	txRecordRepo t.TxRecordRepo
}

func NewService(
	accountRepo account.Repository,
	txRecordRepo t.TxRecordRepo,
) Service {
	return &service{
		accountRepo:  accountRepo,
		txRecordRepo: txRecordRepo,
	}
}

type CreateRequest struct {
	Address     string `json:"address"`
	Hash        string `json:"hash"`
	ProjectName string `json:"project_name"`
	Url         string `json:"url"`
	Amount      string `json:"amount"`
	Symbol      string `json:"symbol"`
	Signature   string `json:"signature"`
}

type CreateResponse struct {
	Succeed bool `json:"succeed"`
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	var err error
	var account *account.Account

	fmt.Println("req.Address", req.Address)

	if account, err = s.accountRepo.FirstOrCreate(ctx, req.Address); err != nil {
		return nil, err
	}

	var tx t.TxRecord
	tx.Account = account.Address
	tx.Hash = req.Hash
	tx.Signature = req.Signature
	if err = s.txRecordRepo.Create(ctx, &tx); err != nil {
		return nil, err
	}

	return &CreateResponse{Succeed: true}, nil
}

type ListRequest struct {
	Address string `json:"address"`
}

type txRecord struct {
	Address     string `gorm:"column:address;type:varchar(42)"`
	Hash        string `gorm:"column:hash;type:varchar(255)"`
	ProjectName string `gorm:"column:project_name;type:varchar(255)"`
	Url         string `gorm:"column:url;type:varchar(2048)"`
	Amount      string `gorm:"column:amount;type:varchar(255)"`
	Symbol      string `gorm:"column:symbol;type:varchar(255)"`
	Signature   string `gorm:"column:signature;type:varchar(255)"`
}

type ListResponse struct {
	TxRecords []txRecord `json:"tx_records"`
}

func (s *service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	var err error
	var txRecords *[]t.TxRecord
	if txRecords, err = s.txRecordRepo.ListByAddress(ctx, req.Address); err != nil {
		return nil, err
	}

	var res ListResponse
	for _, txRecord := range *txRecords {
		res.TxRecords = append(res.TxRecords, txRecordToResponse(txRecord))
	}

	return &res, nil
}

func txRecordToResponse(t t.TxRecord) txRecord {
	return txRecord{
		Address:   t.Account,
		Hash:      t.Hash,
		Signature: t.Signature,
	}
}
