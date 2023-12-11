package service

import (
	"cart-backend/internal/domain/account"
	txrecord "cart-backend/internal/domain/tx_record"
	"context"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*CreateResponse, error)
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
}

type service struct {
	accountRepo  account.Repository
	txRecordRepo txrecord.Repository
}

func NewService(
	accountRepo account.Repository,
	txRecordRepo txrecord.Repository,
) Service {
	return &service{
		accountRepo:  accountRepo,
		txRecordRepo: txRecordRepo,
	}
}

type CreateRequest struct {
	Address     string `json:"address"`
	ProjectName string `json:"project_name"`
	Url         string `json:"url"`
	Amount      string `json:"amount"`
	Symbol      string `json:"symbol"`
}

type CreateResponse struct {
	Succeed bool `json:"succeed"`
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	var err error
	var account *account.Account
	if account, err = s.accountRepo.FirstOrCreate(ctx, req.Address); err != nil {
		return nil, err
	}

	var txRecord txrecord.TxRecord
	txRecord.Account = *account
	txRecord.ProjectName = req.ProjectName
	txRecord.Url = req.Url
	txRecord.Amount = req.Amount
	txRecord.Symbol = req.Symbol
	if err = s.txRecordRepo.Create(ctx, &txRecord); err != nil {
		return nil, err
	}

	return &CreateResponse{Succeed: true}, nil
}

type ListRequest struct {
	Address string `json:"address"`
}

type txRecord struct {
	Address     string `gorm:"column:address;type:varchar(42)" json:"address"`
	ProjectName string `gorm:"column:project_name;type:varchar(255)" json:"project_name"`
	Url         string `gorm:"column:url;type:varchar(2048)" json:"url"`
	Amount      string `gorm:"column:amount;type:varchar(255)" json:"amount"`
	Symbol      string `gorm:"column:symbol;type:varchar(255)" json:"symbol"`
}

type ListResponse struct {
	TxRecords []txRecord `json:"tx_records"`
}

func (s *service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	var err error
	var txRecords *[]txrecord.TxRecord
	if txRecords, err = s.txRecordRepo.ListByAddress(ctx, req.Address); err != nil {
		return nil, err
	}

	var res ListResponse
	for _, txRecord := range *txRecords {
		res.TxRecords = append(res.TxRecords, txRecordToResponse(txRecord))
	}

	return &res, nil
}

func txRecordToResponse(t txrecord.TxRecord) txRecord {
	return txRecord{
		Address:     t.Account.Address,
		ProjectName: t.ProjectName,
		Url:         t.Url,
		Amount:      t.Amount,
		Symbol:      t.Symbol,
	}
}
