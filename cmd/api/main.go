package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cart-backend/internal/adapter"
	"cart-backend/internal/app/api"
	"cart-backend/internal/service"
	"cart-backend/pkg/app"
	"cart-backend/pkg/log"

	gormpkg "cart-backend/pkg/gorm"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// init logger
	sync, err := log.Init(log.Config{
		Name:   "cart-backend.api",
		Level:  zapcore.InfoLevel,
		Stdout: true,
		File:   "log/cart-backend/api.log",
	})
	if err != nil {
		panic(err)
	}
	defer sync()
	logger := zap.L()

	// prepare context
	ctx := app.GraceCtx(context.Background())

	// init db
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := gormpkg.NewGormPostgresConn(
		gormpkg.Config{
			DSN:             dsn,
			MaxIdleConns:    2,
			MaxOpenConns:    2,
			ConnMaxLifetime: 10 * time.Minute,
			SingularTable:   true,
		},
	)
	if err != nil {
		logger.Error("connect to database error", zap.Error(err))
		return
	}

	// prepare service
	accountRepo := adapter.NewAccountRepo(db)
	txRecordRepo := adapter.NewTxRecordRepo(db)

	svc := service.NewService(
		accountRepo,
		txRecordRepo,
	)

	app := api.New(api.Config{
		Port:    80,
		Service: svc,
	})
	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("app.Start", zap.Error(err))
	}
}
