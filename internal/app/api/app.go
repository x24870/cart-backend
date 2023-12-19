package api

import (
	"cart-backend/internal/handler"
	"cart-backend/pkg/api/middlewares"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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

	// CORS
	config := cors.DefaultConfig()
	// config.AllowAllOrigins = true // adjust this to your needs
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowOrigins = []string{"https://innovatechain.xyz", "https://api.innovatechain.xyz", "http://localhost:5173", "http://localhost:8080"}
	router.Use(cors.New(config))

	// router.Use(middlewares.PaniwcCatcher)
	router.Use(middlewares.CtxLogger)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/tx_record", a.Handler.CreateTxRecord)
	router.POST("/tx_record/list", a.Handler.ListTxRecordByAddress)

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
	if os.Getenv("ENV") == "dev" {
		go server.ListenAndServe()
	} else {
		go func() {
			// Serve HTTP, which will redirect to HTTPS
			h := m.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		// Start HTTPS server
		go func() {
			log.Fatal(server.ListenAndServeTLS("", "")) // Key and cert are coming from Let's Encrypt
		}()
	}

	<-ctx.Done()
	return server.Shutdown(ctx)
}
