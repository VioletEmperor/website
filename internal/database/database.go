package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return pool, nil
}
