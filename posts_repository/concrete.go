package posts_repository

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

type ConcretePostsRepository struct {
    Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) ConcretePostsRepository {
    return ConcretePostsRepository{pool}
}

func (repo ConcretePostsRepository) GetPost(id int) (*Post, error) {
    query := "SELECT * FROM public.posts WHERE id = $1"

    row := repo.Pool.QueryRow(context.Background(), query, id)

    var post Post

    err := row.Scan(
        &post.ID,
        &post.Title,
        &post.Author,
        &post.Created,
        &post.Edited,
        &post.Body,
    )

    if err != nil {
        return nil, fmt.Errorf("error scanning post: %v", err)
    }

    return &post, nil
}

func (repo ConcretePostsRepository) GetPosts() ([]Post, error) {
    var posts []Post

    query := "SELECT * FROM public.posts"

    rows, err := repo.Pool.Query(context.Background(), query)

    if err != nil {
        return nil, fmt.Errorf("error getting posts: %w", err)
    }

    for rows.Next() {
        var post Post

        err := rows.Scan(
            &post.ID,
            &post.Title,
            &post.Author,
            &post.Created,
            &post.Edited,
            &post.Body,
        )

        if err != nil {
            return nil, fmt.Errorf("error scanning post: %w", err)
        }

        posts = append(posts, post)
    }

    return posts, nil
}
