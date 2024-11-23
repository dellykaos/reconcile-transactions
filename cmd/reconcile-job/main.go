package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/delly/amartha/config"
	filestorage "github.com/delly/amartha/repository/file_storage"
	localfilestorage "github.com/delly/amartha/repository/file_storage/local_file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
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

	querier := dbgen.New(pool)
	var fileStorage filestorage.FileStorageRepository
	if cfg.LocalStorage.UseLocal {
		currentDir, err := os.Getwd()
		checkError(err)
		fileDir := fmt.Sprintf("%s%s", currentDir, cfg.LocalStorage.Dir)
		fileStorage = localfilestorage.NewStorage(fileDir)
	}
	reconProcesserService := reconciliatonjob.NewProcesserService(querier, fileStorage)

	logger.Println("Processing reconciliation job...")
	err = reconProcesserService.Process(ctx)
	checkError(err)
	logger.Println("Reconciliation job processed")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
