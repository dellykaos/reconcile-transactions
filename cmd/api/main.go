package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/delly/amartha/common"
	"github.com/delly/amartha/config"
	handler "github.com/delly/amartha/handler/http"
	filestorage "github.com/delly/amartha/repository/file_storage"
	localfilestorage "github.com/delly/amartha/repository/file_storage/local_file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig(".env")
	checkError(err)

	common.SetupLogger(cfg.Env)
	logger := common.Logger()
	logger.Info("Connecting to database...")
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
	reconFinderSvc := reconciliatonjob.NewFinderService(querier)
	reconCreatorSvc := reconciliatonjob.NewCreatorService(querier, fileStorage)
	reconJobHandler := handler.NewReconciliationJobHandler(reconFinderSvc, reconCreatorSvc)

	r := httprouter.New()
	reconJobHandler.Register(r)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		logger.Info(fmt.Sprintf("server is listening at %s", s.Addr))
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Panic(err.Error(), zapcore.Field{
					Key:       "error",
					Interface: err,
				})
			}
		}
	}(&srv)
	<-signalChan

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("server shutdown failed: %v", err))
		return
	}

	logger.Info("exiting...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
