package main

import (
	"context"
	"fmt"

	"github.com/delly/amartha/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig(".env")
	checkError(err)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	pool, err := pgxpool.Connect(ctx, connStr)
	checkError(err)
	defer pool.Close()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
