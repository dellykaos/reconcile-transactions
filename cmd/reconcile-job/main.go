package main

import (
	"context"
	"fmt"
	"os"

	"github.com/delly/amartha/common"
	"github.com/delly/amartha/config"
	filestorage "github.com/delly/amartha/repository/file_storage"
	localfilestorage "github.com/delly/amartha/repository/file_storage/local_file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig(".env")
	checkError(err)

	common.SetupLogger(cfg.Env)
	logger := common.Logger()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	pool, err := pgxpool.Connect(ctx, connStr)
	checkError(err)
	logger.Info("Connected to database")
	defer func() {
		logger.Info("Closing database connection...")
		pool.Close()
		logger.Info("Database connection closed")
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

	logger.Info("Processing reconciliation job...")
	err = reconProcesserService.Process(ctx)
	checkError(err)
	logger.Info("Reconciliation job processed")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
