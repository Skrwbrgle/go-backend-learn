package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Skrwbrgle/go-backend-learn/config"
	dbpkg "github.com/Skrwbrgle/go-backend-learn/internal/db"
	"github.com/Skrwbrgle/go-backend-learn/internal/middleware"
	"github.com/Skrwbrgle/go-backend-learn/internal/routes"
	"github.com/Skrwbrgle/go-backend-learn/pkg/logger"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp() *App {
	cfg := config.LoadConfig()

	// Init logger
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	logger.InitLogger(env)
	defer func() {
		if err := logger.Log.Sync(); err != nil {
			logger.Log.Error("failed to sync logger", zap.Error(err))
		}
	}()
	logger.Log.Info("Starting server...", zap.String("env", env))

	// DB
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	dbpkg.AutoMigrate(db)

	r := gin.Default()

	// Middleware stack (order matters!)
	r.Use(middleware.CorrelationIDMiddleware())     // 1. Add correlation ID ke setiap request
	r.Use(middleware.LoggerMiddleware(logger.Log))  // 2. Log request
	r.Use(middleware.MetricsMiddleware())           // 3. Record metrics

	// Router
	routes.SetupRouter(r, db)

	return &App{
		Router: r,
		DB:     db,
	}
}

// Run starts the server with graceful shutdown
func (a *App) Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	// Jalankan server di goroutine
	go func() {
		log.Println("🚀 Server running on", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %s\n", err)
		}
	}()

	// Tangkap signal interrupt/terminate
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("⚡ Shutting down server...")

	// Shutdown dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %s", err)
	}

	// Tutup DB pool
	if sqlDB, err := a.DB.DB(); err == nil {
		_ = sqlDB.Close()
	}

	log.Println("✅ Server exited gracefully")
}
