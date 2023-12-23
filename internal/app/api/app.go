package api

import (
	"cart-backend/internal/handler"
	"cart-backend/pkg/api/middlewares"
	"context"
	"crypto/tls"
	"log"
	"net/http"

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
	config.AllowOrigins = []string{
		"https://innovatechain.xyz",
		"http://innovatechain.xyz",
		"https://api.innovatechain.xyz",
		"http://api.innovatechain.xyz",
		"https://api2.innovatechain.xyz",
		"http://api2.innovatechain.xyz",
		"https://sharing.innovatechain.xyz",
		"http://sharing.innovatechain.xyz",
		"http://localhost:5173",
		"http://localhost:8080",
	}
	router.Use(cors.New(config))

	// router.Use(middlewares.PaniwcCatcher)
	router.Use(middlewares.CtxLogger)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/tx_record", a.Handler.CreateTxRecord)
	router.POST("/tx_record/list", a.Handler.ListTxRecordByAddress)

	// Set up the HTTP server
	httpServer := &http.Server{
		Addr:    "0.0.0.0:80", // HTTP port
		Handler: router,
	}

	// Let's Encrypt tls certificate
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.innovatechain.xyz", "api2.innovatechain.xyz"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	// Set up the HTTPS server
	httpsServer := &http.Server{
		Addr:      "0.0.0.0:443", // HTTPS port
		Handler:   router,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}
	httpsServer.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	// Start Running Server.
	serveBothHttpAndHttps := false
	if serveBothHttpAndHttps {
		// Start the HTTP server in a new goroutine
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe error: %v", err)
			}
		}()

		// Start the HTTPS server in a new goroutine
		go func() {
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ListenAndServeTLS error: %v", err)
			}
		}()
	} else {
		go func() {
			// Serve HTTP, which will redirect to HTTPS
			h := m.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		// Start HTTPS server
		go func() {
			log.Fatal(httpsServer.ListenAndServeTLS("", "")) // Key and cert are coming from Let's Encrypt
		}()
	}

	<-ctx.Done()
	// Shutdown both HTTP and HTTPS servers
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("HTTP Server Shutdown:", err)
	}
	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Fatal("HTTPS Server Shutdown:", err)
	}

	return nil
}
