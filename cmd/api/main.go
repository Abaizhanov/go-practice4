package main

import (
	"context"
	"log"
	"time"

	"practice4-sqlx/internal/config"
	"practice4-sqlx/internal/db"
	"practice4-sqlx/internal/repository/postgres"
	"practice4-sqlx/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	database, err := db.Open(cfg.DSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer database.Close()

	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := postgres.NewUserRepo()
	svc := service.NewUserService(database, repo)

	ctx := context.Background()

	_ = svc.InsertUser(ctx, "Alice", "alice@example.com", 100)
	_ = svc.InsertUser(ctx, "Bob", "bob@example.com", 50)

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	if err := svc.TransferBalance(ctx, 1, 2, 25.5); err != nil {
		log.Printf("transfer failed: %v", err)
	} else {
		log.Println("transfer ok: 1 -> 2 amount=25.5")
	}
}
