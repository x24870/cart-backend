package service

import (
	"cart-backend/internal/domain/account"
	t "cart-backend/internal/domain/transaction"
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*CreateResponse, error)
	List(ctx context.Context, req ListRequest) (*ListResponse, error)
}

type service struct {
	accountRepo   account.Repository
	txRecordRepo  t.TxRecordRepo
	operationRepo t.OperationRepo
	intentRepo    t.IntentRepo
}

func NewService(
	accountRepo account.Repository,
	txRecordRepo t.TxRecordRepo,
	operationRepo t.OperationRepo,
	intentRepo t.IntentRepo,
) Service {
	return &service{
		accountRepo:   accountRepo,
		txRecordRepo:  txRecordRepo,
		operationRepo: operationRepo,
		intentRepo:    intentRepo,
	}
}

type svcTxRecord struct {
	Address    string         `json:"address"`
	Hash       string         `json:"hash"`
	Signature  string         `json:"signature"`
	Operations []svcOperation `json:"operations"`
}

type svcOperation struct {
	ProjectName string      `json:"project_name"`
	Url         string      `json:"url"`
	Symbol      string      `json:"symbol"`
	Amount      string      `json:"amount"`
	Decimal     uint        `json:"decimal"`
	Intents     []svcIntent `json:"intents"`
}

type svcIntent struct {
	Description string `json:"description"`
}

type CreateRequest struct {
	svcTxRecord
}

func createReqToTxRecord(req CreateRequest) t.TxRecord {
	var txRecord t.TxRecord
	txRecord.Account = req.Address
	txRecord.Hash = req.Hash
	txRecord.Signature = req.Signature
	return txRecord
}

func createReqToOperations(req CreateRequest) ([]t.Operation, []t.Intent) {
	var operations []t.Operation
	var intents []t.Intent
	for _, operation := range req.Operations {
		var op t.Operation
		op.ID = uuid.NewV4()
		op.TxHash = req.Hash
		op.ProjectName = operation.ProjectName
		op.Url = operation.Url
		op.Amount = operation.Amount
		op.Symbol = operation.Symbol
		op.Decimal = operation.Decimal
		operations = append(operations, op)

		for _, intent := range operation.Intents {
			var it t.Intent
			it.OperationID = op.ID
			it.Description = intent.Description
			intents = append(intents, it)
		}
	}
	return operations, intents
}

type CreateResponse struct {
	Succeed bool `json:"succeed"`
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	if _, err := s.accountRepo.FirstOrCreate(ctx, req.Address); err != nil {
		return nil, err
	}

	tx := createReqToTxRecord(req)
	ops, its := createReqToOperations(req)
	if err := s.txRecordRepo.Create(ctx, &tx, &ops, &its); err != nil {
		return nil, err
	}

	return &CreateResponse{Succeed: true}, nil
}

type ListRequest struct {
	Address string `json:"address"`
}

type ListResponse struct {
	TxRecords []svcTxRecord `json:"tx_records"`
}

func (s *service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	var err error
	var txRecords *[]t.TxRecord
	var svcTxRecords []svcTxRecord
	if txRecords, err = s.txRecordRepo.ListByAddress(ctx, req.Address); err != nil {
		return nil, err
	}
	for _, txRecord := range *txRecords {
		svcTxRecords = append(svcTxRecords, txRecordToResponse(txRecord))
	}
	fmt.Printf("svcTxRecords: %+v\n", svcTxRecords)

	var res ListResponse
	for _, svcTx := range svcTxRecords {
		var operations *[]t.Operation
		if operations, err = s.operationRepo.ListByTxHash(ctx, svcTx.Hash); err != nil {
			return nil, err
		}

		for _, operation := range *operations {
			var svcOp svcOperation
			svcOp.ProjectName = operation.ProjectName
			svcOp.Url = operation.Url
			svcOp.Amount = operation.Amount
			svcOp.Symbol = operation.Symbol
			svcOp.Decimal = operation.Decimal
			fmt.Printf("svcOp: %+v\n", svcOp)

			var intents *[]t.Intent
			if intents, err = s.intentRepo.ListByOperationID(ctx, operation.ID); err != nil {
				return nil, err
			}

			for _, intent := range *intents {
				svcOp.Intents = append(svcOp.Intents, svcIntent{Description: intent.Description})
				fmt.Printf("intent: %+v\n", intent)
			}

			svcTx.Operations = append(svcTx.Operations, svcOp)
		}

		res.TxRecords = append(res.TxRecords, svcTx)
	}

	return &res, nil
}

func txRecordToResponse(t t.TxRecord) svcTxRecord {
	return svcTxRecord{
		Address:   t.Account,
		Hash:      t.Hash,
		Signature: t.Signature,
	}
}
