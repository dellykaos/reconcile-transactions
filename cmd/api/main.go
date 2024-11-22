package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/delly/amartha/config"
	handler "github.com/delly/amartha/handler/http"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
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
	reconFinderSvc := reconciliatonjob.NewFinderService(querier)
	reconJobHandler := handler.NewReconciliationJobHandler(*reconFinderSvc)

	r := httprouter.New()
	reconJobHandler.Register(r)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		log.Printf("server is listening at %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Panic(err)
			}
		}
	}(&srv)
	<-signalChan

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
		return
	}

	logger.Println("exiting...")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}