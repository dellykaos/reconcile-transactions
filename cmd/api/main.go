package main

import (
	"context"
	"fmt"
	"log"

	"github.com/delly/amartha/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	logger := log.Default()
	logger.Println("Loading configuration...")
	cfg, err := config.NewConfig(".env")
	checkError(err)
	logger.Println("Configuration loaded")

	logger.Println("Connecting to database...")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	pool, err := pgxpool.Connect(ctx, connStr)
	checkError(err)
	logger.Println("Connected to database")
	defer func() {
		logger.Println("Closing database connection...")
		pool.Close()
		logger.Println("Database connection closed")
	}()

	// setup http server

	logger.Println("exiting...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
