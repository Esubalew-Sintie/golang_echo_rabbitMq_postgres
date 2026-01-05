package initiator

import (
	"context"
	"net/http"
	"payment-gateway/internal/handler/middleware"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

func Initialize(ctx context.Context) {
	log := InitLogger(ctx)

	log.Info(ctx, "initializing config")
	cfg, err := InitConfig(ctx, log)
	if err != nil {
		log.Fatal(ctx, "failed to initialize config: %v", err)
	}
	log.Info(ctx, "config initialized")

	pool, err := InitPostgres(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialize database: %v", err)
	}
	defer ClosePostgres(pool)

	log.Info(ctx, "database initialized")

	persistence := InitPersistence(pool, log)
	log.Info(ctx, "persistence initialized")

	messaging, err := InitMessageBroker(log)
	if err != nil {
		log.Fatal(ctx, "failed to initialize message broker: %v", err)
	}
	defer func() {
		if err := messaging.Close(); err != nil {
			log.Error(ctx, "Failed to close messaging: %v", err)
		}
	}()
	log.Info(ctx, "messaging initialized")

	svc := InitService(persistence, messaging, log)
	log.Info(ctx, "services initialized")

	hdlr := InitHandler(svc, log)
	log.Info(ctx, "handlers initialized")

	authMiddleware := middleware.NewAuthMiddleware(cfg.JwtSecretKey)
	log.Info(ctx, "auth middleware initialized")

	e := echo.New()
	e.Use(middleware.EchoRecovery())
	e.Use(middleware.EchoLogger())
	e.Use(middleware.EchoCORS())
	e.Use(middleware.EchoSecurityHeaders())
	e.Use(middleware.EchoTimeout(30 * time.Second))

	InitRoutes(e, hdlr, authMiddleware, log)
	log.Info(ctx, "router initialized")

	worker := InitWorker(svc, messaging, log)
	log.Info(ctx, "worker initialized")

	serverAddr := cfg.ServerAddr
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	go func() {
		log.Info(ctx, "starting HTTP server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, "HTTP server failed: %v", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.Start(ctx)
	}()

	log.Info(ctx, "payment gateway service started successfully")

	<-ctx.Done()
	log.Info(ctx, "shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info(shutdownCtx, "shutting down HTTP server")
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Error(shutdownCtx, "HTTP server shutdown error: %v", err)
	}

	log.Info(shutdownCtx, "waiting for worker to finish")
	wg.Wait()

	log.Info(shutdownCtx, "service shutdown completed")
}
