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
	query := "SELECT * FROM public.posts ORDER BY created DESC"

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

func (repo ConcreteRepository) GetPostsPaginated(page int) ([]Post, PaginationInfo, error) {
	// Get total count first
	totalPosts, err := repo.GetTotalPostsCount()
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("error getting total posts count: %w", err)
	}

	// Create pagination info
	paginationInfo := NewPaginationInfo(totalPosts, page)
	
	// Get paginated posts
	query := "SELECT * FROM public.posts ORDER BY created DESC LIMIT $1 OFFSET $2"
	
	rows, err := repo.Pool.Query(context.Background(), query, PostsPerPage, paginationInfo.GetOffset())

	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("error getting paginated posts: %w", err)
	}

	posts, err := pgx.CollectRows[Post](rows, pgx.RowToStructByName[Post])

	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("error scanning paginated posts: %w", err)
	}

	return posts, paginationInfo, nil
}

func (repo ConcreteRepository) GetTotalPostsCount() (int, error) {
	query := "SELECT COUNT(*) FROM public.posts"
	
	var count int
	err := repo.Pool.QueryRow(context.Background(), query).Scan(&count)
	
	if err != nil {
		return 0, fmt.Errorf("error getting posts count: %w", err)
	}
	
	return count, nil
}

func (repo ConcreteRepository) DeletePost(id int) error {
	query := "DELETE FROM public.posts WHERE id = $1"
	
	result, err := repo.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("post with id %d not found", id)
	}
	
	return nil
}

func (repo ConcreteRepository) UpdatePost(id int, title, description, body string) error {
	query := `UPDATE public.posts 
		SET title = $2, description = $3, body = $4, edited = NOW() 
		WHERE id = $1`
	
	result, err := repo.Pool.Exec(context.Background(), query, id, title, description, body)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("post with id %d not found", id)
	}
	
	return nil
}

func (repo ConcreteRepository) CreatePost(title, description, body, author string) error {
	query := `INSERT INTO public.posts (title, description, body, author, created, edited) 
		VALUES ($1, $2, $3, $4, NOW(), NOW())`
	
	_, err := repo.Pool.Exec(context.Background(), query, title, description, body, author)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}
	
	return nil
}
