package api

import (
	"cart-backend/internal/handler"
	"cart-backend/pkg/api/middlewares"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
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
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890","project_name":"rysk","url":"https://rysk.fi","amount":"10","symbol":"USDC","signature":"0x21fbf0696d5e0aa2ef41a2b4ffb623bcaf070461d61cf7251c74161f82fec3a4370854bc0a34b3ab487c1bc021cd318c734c51ae29374f2beb0e6f2dd49b4bf41c"}' http://localhost:8080/tx_record
	// curl -X POST -H "Content-Type: application/json" -d '{"address":"0x1234567890123456789012345678901234567890"}' http://localhost:8080/tx_record/list

	// Setup HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", a.Port),
		Handler: router,
	}

	// Let's Encrypt tls certificate
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.innovatechain.xyz"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}
	server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	// Start Running HTTP Server.
	// go server.ListenAndServe()

	go func() {
		// Serve HTTP, which will redirect to HTTPS
		h := m.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	// Start HTTPS server
	go func() {
		log.Fatal(server.ListenAndServeTLS("", "")) // Key and cert are coming from Let's Encrypt
	}()

	<-ctx.Done()
	return server.Shutdown(ctx)
}
