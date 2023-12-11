package service

import (
	"cart-backend/internal/domain/account"
	txrecord "cart-backend/internal/domain/tx_record"
	"context"
)

type Service interface {
	List(ctx context.Context, address string) ([]account.Account, error)
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

func (s *service) List(ctx context.Context, address string) ([]account.Account, error) {
	// return s.repo.List(ctx)
	return nil, nil
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
