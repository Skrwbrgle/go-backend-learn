package dbpkg

import (
	"fmt"
	"log"

	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	// Import file source for database migration
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Import postgres driver
	_ "github.com/lib/pq"
)

func RunMigrations(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error running migration: %v", err)
	}

	fmt.Println("✅ Database migration completed!")
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		// &model.Product{},
		// &model.Order{},
	)
	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}
}
