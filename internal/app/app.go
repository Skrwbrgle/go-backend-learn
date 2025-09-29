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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Skrwbrgle/go-backend-learn/config"
	dbpkg "github.com/Skrwbrgle/go-backend-learn/internal/db"
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
	logger.InitLogger()
	defer logger.Log.Sync()

	// DB
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	dbpkg.AutoMigrate(db)

	// Router
	r := gin.Default()
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