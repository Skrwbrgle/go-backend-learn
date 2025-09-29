package main

import (
	"log"

	"github.com/Skrwbrgle/go-backend-learn/config"
	dbpkg "github.com/Skrwbrgle/go-backend-learn/internal/db"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	dsn := cfg.DBUrl

	dbConn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbConn.Close()

	dbpkg.RunMigrations(dbConn)
}
