package api

import (
	"cart-backend/internal/service"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port    int
	Service service.Service
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

	// router.Use(middlewares.PanicCatcher)
	// router.Use(cors.Default())
	// router.Use(middlewares.CtxLogger)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

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
