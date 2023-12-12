package api

import (
	"cart-backend/internal/handler"
	"cart-backend/pkg/api/middlewares"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port int
	// Service service.Service
	Handler handler.Handler
}

func New(config Config) *App {
	return &App{config}
}

type App struct {
	Config
}

func (a *App) Start(ctx context.Context) error {
	// start a gin HTTP server
	// Create our HTTP Router
	router := gin.New()

	// Configure HTTP Router Settings
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = false
	router.HandleMethodNotAllowed = false
	router.ForwardedByClientIP = true
	router.AppEngine = false
	router.UseRawPath = false
	router.UnescapePathValues = true
	router.ContextWithFallback = true

	// router.Use(middlewares.PaniwcCatcher)
	// router.Use(cors.Default())
	router.Use(middlewares.CtxLogger)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/tx_record", a.Handler.CreateTxRecord)
	router.POST("/tx_record/list", a.Handler.ListTxRecordByAddress)

	// curl example:
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890","project_name":"rysk","url":"https://rysk.fi","amount":"10","symbol":"USDC"}' http://localhost:8080/tx_record
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890","project_name":"rysk","url":"https://rysk.fi","amount":"10","symbol":"USDC"}' http://localhost:8080/tx_record | jq
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890"}' http://localhost:8080/tx_record/list
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890"}' http://localhost:8080/tx_record/list | jq

	// Setup HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.Port),
		Handler: router,
	}

	// Start Running HTTP Server.
	go server.ListenAndServe()
	<-ctx.Done()
	return server.Shutdown(ctx)
}
