package service

import (
	"cart-backend/internal/domain/account"
	t "cart-backend/internal/domain/transaction"
	"context"
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
	Address    string `json:"address"`
	Hash       string `json:"hash"`
	Signature  string `json:"signature"`
	Operations []struct {
		ProjectName string `json:"project_name"`
		Url         string `json:"url"`
		Amount      string `json:"amount"`
		Symbol      string `json:"symbol"`
		Intents     []struct {
			Description string `json:"description"`
		} `json:"intents"`
	} `json:"operations"`
}

func createReqToTxRecord(req CreateRequest) t.TxRecord {
	var txRecord t.TxRecord
	txRecord.Account = req.Address
	txRecord.Hash = req.Hash
	txRecord.Signature = req.Signature
	return txRecord
}

func createReqToOperations(req CreateRequest) []t.Operation {
	var operations []t.Operation
	for _, operation := range req.Operations {
		var op t.Operation
		op.ProjectName = operation.ProjectName
		op.Url = operation.Url
		op.Amount = operation.Amount
		op.Symbol = operation.Symbol
		operations = append(operations, op)
	}
	return operations
}

func createReqToIntents(req CreateRequest) []t.Intent {
	var intents []t.Intent
	for _, operation := range req.Operations {
		for _, intent := range operation.Intents {
			var it t.Intent
			it.Description = intent.Description
			intents = append(intents, it)
		}
	}
	return intents
}

type CreateResponse struct {
	Succeed bool `json:"succeed"`
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if _, err := s.accountRepo.FirstOrCreate(ctx, req.Address); err != nil {
		return nil, err
	}

	tx := createReqToTxRecord(req)
	ops := createReqToOperations(req)
	its := createReqToIntents(req)
	if err := s.txRecordRepo.Create(ctx, &tx, &ops, &its); err != nil {
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
