package main

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

func connect(url string) (*pgxpool.Pool, error) {
    pool, err := pgxpool.New(context.Background(), url)

    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %v", err)
    }

    return pool, nil
}
