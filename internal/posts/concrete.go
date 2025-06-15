package posts

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type ConcreteRepository struct {
    Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) ConcreteRepository {
    return ConcreteRepository{pool}
}

func (repo ConcreteRepository) GetPost(id int) (*Post, error) {
    query := "SELECT * FROM public.posts WHERE id = $1"

    row, err := repo.Pool.Query(context.Background(), query, id)

    if err != nil {
        return nil, fmt.Errorf("error getting posts: %w", err)
    }

    post, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Post])

    if err != nil {
        return nil, err
    }

    return &post, nil
}

func (repo ConcreteRepository) GetPosts() ([]Post, error) {
    query := "SELECT * FROM public.posts"

    rows, err := repo.Pool.Query(context.Background(), query)

    if err != nil {
        return nil, fmt.Errorf("error getting posts: %w", err)
    }

    posts, err := pgx.CollectRows[Post](rows, pgx.RowToStructByName[Post])

    if err != nil {
        return nil, fmt.Errorf("error scanning posts: %w", err)
    }

    return posts, nil
}
