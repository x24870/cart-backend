package service

import (
	"cart-backend/internal/domain/account"
	txrecord "cart-backend/internal/domain/tx_record"
	"context"
)

type Service interface {
	Create(ctx context.Context, txRecord *txrecord.TxRecord) error
	List(ctx context.Context, address string) (*[]txrecord.TxRecord, error)
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

type CreateTxRecordReq struct {
	Address     string `json:"address"`
	ProjectName string `json:"project_name"`
}

func (s *service) Create(ctx context.Context, txRecord *txrecord.TxRecord) error {
	var err error
	if _, err = s.accountRepo.FirstOrCreate(ctx, *txRecord.Address); err != nil {
		return err
	}
	if err := s.txRecordRepo.Create(ctx, txRecord); err != nil {
		return err
	}

	return nil
}

func (s *service) List(ctx context.Context, address string) (*[]txrecord.TxRecord, error) {
	// return s.repo.List(ctx)
	return s.txRecordRepo.ListByAddress(ctx, address)
}

// func (s *service) Execute(ctx context.Context, lambdaName string, input []byte) ([]byte, error) {
// 	lambda, err := s.repo.GetLambdaByName(ctx, lambdaName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	runtime, ok := s.runtimes[lambda.Metadata.RuntimeType]
// 	if !ok {
// 		return nil, fmt.Errorf(
// 			"runtime %s not found; %w",
// 			lambda.Metadata.RuntimeType, domain.ErrRuntimeNotFound,
// 		)
// 	}
// 	return runtime.Exec(ctx, lambda.Code, input)
// }
