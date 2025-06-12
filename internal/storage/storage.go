package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB  *pgxpool.Pool
	Ctx context.Context
}

func NewPostgres(ctx context.Context, dbURI string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize db storage: %v", err)
	}

	return &Postgres{
		DB:  pool,
		Ctx: ctx,
	}, nil
}


