package main

import (
	"context"

	"cart-backend/internal/app/api"
	"cart-backend/pkg/app"
)

func main() {
	// init logger
	// sync, err := log.Init(log.Config{
	// 	Name:   "lambda-protocol.api",
	// 	Level:  zapcore.InfoLevel,
	// 	Stdout: true,
	// 	File:   "log/lambda-protocol/api.log",
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer sync()
	// logger := zap.L()

	// prepare context
	ctx := app.GraceCtx(context.Background())

	// prepare service
	// runtimeV8 := adapter.NewRuntimeV8()
	// // repo := adapter.NewFSRepository("./lambda")
	// repo := adapter.NewIPFSRepository(
	// 	"http://localhost:5001",
	// 	"QmRQgLcZAzbEZN6mhwEySmw2baqFAh1tsoMFxvia82ENPd",
	// )
	// svc := service.NewService(
	// 	repo,
	// 	map[domain.RuntimeType]domain.Runtime{
	// 		domain.RuntimeTypeV8: runtimeV8,
	// 	},
	// )

	app := api.New(api.Config{
		Port: 9999,
		// Service: svc,
	})
	err := app.Start(ctx)
	if err != nil {
		// logger.Fatal("app.Start", zap.Error(err))
	}
}
