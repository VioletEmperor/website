package posts

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestConcreteRepository_GetPost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful get post", func(t *testing.T) {
		expectedPost := Post{
			ID:          1,
			Title:       "Test Post",
			Author:      "Adam Shkolnik",
			Created:     time.Now(),
			Edited:      time.Now(),
			Body:        "test-content.html",
			Description: "Test description",
		}

		mock.ExpectQuery(`SELECT \* FROM public\.posts WHERE id = \$1`).
			WithArgs(1).
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author", "created", "edited", "body", "description"}).
				AddRow(expectedPost.ID, expectedPost.Title, expectedPost.Author, expectedPost.Created, expectedPost.Edited, expectedPost.Body, expectedPost.Description))

		post, err := repo.GetPost(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if post.ID != expectedPost.ID {
			t.Errorf("expected ID %d, got %d", expectedPost.ID, post.ID)
		}
		if post.Title != expectedPost.Title {
			t.Errorf("expected Title %s, got %s", expectedPost.Title, post.Title)
		}
		if post.Author != expectedPost.Author {
			t.Errorf("expected Author %s, got %s", expectedPost.Author, post.Author)
		}
		if post.Body != expectedPost.Body {
			t.Errorf("expected Body %s, got %s", expectedPost.Body, post.Body)
		}
		if post.Description != expectedPost.Description {
			t.Errorf("expected Description %s, got %s", expectedPost.Description, post.Description)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM public\.posts WHERE id = \$1`).
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		post, err := repo.GetPost(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if post != nil {
			t.Error("expected nil post, got non-nil")
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			t.Errorf("expected ErrNoRows, got %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM public\.posts WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(pgx.ErrTxClosed)

		post, err := repo.GetPost(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if post != nil {
			t.Error("expected nil post, got non-nil")
		}
		if err != nil && !contains(err.Error(), "error getting posts") {
			t.Errorf("expected error to contain 'error getting posts', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_GetPosts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful get posts", func(t *testing.T) {
		now := time.Now()
		expectedPosts := []Post{
			{
				ID:          2,
				Title:       "Second Post",
				Author:      "Adam Shkolnik",
				Created:     now,
				Edited:      now,
				Body:        "second-post.html",
				Description: "Second post description",
			},
			{
				ID:          1,
				Title:       "First Post",
				Author:      "Adam Shkolnik",
				Created:     now.Add(-time.Hour),
				Edited:      now.Add(-time.Hour),
				Body:        "first-post.html",
				Description: "First post description",
			},
		}

		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC`).
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author", "created", "edited", "body", "description"}).
				AddRow(expectedPosts[0].ID, expectedPosts[0].Title, expectedPosts[0].Author, expectedPosts[0].Created, expectedPosts[0].Edited, expectedPosts[0].Body, expectedPosts[0].Description).
				AddRow(expectedPosts[1].ID, expectedPosts[1].Title, expectedPosts[1].Author, expectedPosts[1].Created, expectedPosts[1].Edited, expectedPosts[1].Body, expectedPosts[1].Description))

		posts, err := repo.GetPosts()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("expected 2 posts, got %d", len(posts))
		}
		if len(posts) > 0 && posts[0].Title != expectedPosts[0].Title {
			t.Errorf("expected first post title %s, got %s", expectedPosts[0].Title, posts[0].Title)
		}
		if len(posts) > 1 && posts[1].Title != expectedPosts[1].Title {
			t.Errorf("expected second post title %s, got %s", expectedPosts[1].Title, posts[1].Title)
		}
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC`).
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author", "created", "edited", "body", "description"}))

		posts, err := repo.GetPosts()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(posts) != 0 {
			t.Errorf("expected empty posts, got %d", len(posts))
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC`).
			WillReturnError(pgx.ErrTxClosed)

		posts, err := repo.GetPosts()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if posts != nil {
			t.Error("expected nil posts, got non-nil")
		}
		if err != nil && !contains(err.Error(), "error getting posts") {
			t.Errorf("expected error to contain 'error getting posts', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_GetTotalPostsCount(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful count", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(10))

		count, err := repo.GetTotalPostsCount()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 10 {
			t.Errorf("expected count 10, got %d", count)
		}
	})

	t.Run("zero count", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(0))

		count, err := repo.GetTotalPostsCount()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnError(pgx.ErrTxClosed)

		count, err := repo.GetTotalPostsCount()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
		if err != nil && !contains(err.Error(), "error getting posts count") {
			t.Errorf("expected error to contain 'error getting posts count', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_GetPostsPaginated(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful pagination - first page", func(t *testing.T) {
		now := time.Now()
		expectedPosts := []Post{
			{ID: 5, Title: "Post 5", Author: "Adam Shkolnik", Created: now, Edited: now, Body: "post5.html", Description: "Post 5 desc"},
			{ID: 4, Title: "Post 4", Author: "Adam Shkolnik", Created: now.Add(-time.Hour), Edited: now.Add(-time.Hour), Body: "post4.html", Description: "Post 4 desc"},
		}

		// Mock count query
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(12))

		// Mock paginated query
		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(PostsPerPage, 0).
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author", "created", "edited", "body", "description"}).
				AddRow(expectedPosts[0].ID, expectedPosts[0].Title, expectedPosts[0].Author, expectedPosts[0].Created, expectedPosts[0].Edited, expectedPosts[0].Body, expectedPosts[0].Description).
				AddRow(expectedPosts[1].ID, expectedPosts[1].Title, expectedPosts[1].Author, expectedPosts[1].Created, expectedPosts[1].Edited, expectedPosts[1].Body, expectedPosts[1].Description))

		posts, pagination, err := repo.GetPostsPaginated(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(posts) != 2 {
			t.Errorf("expected 2 posts, got %d", len(posts))
		}
		if pagination.CurrentPage != 1 {
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 3 {
			t.Errorf("expected TotalPages 3, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 12 {
			t.Errorf("expected TotalPosts 12, got %d", pagination.TotalPosts)
		}
		if !pagination.HasNext {
			t.Error("expected HasNext to be true")
		}
		if pagination.HasPrev {
			t.Error("expected HasPrev to be false")
		}
		if pagination.NextPage != 2 {
			t.Errorf("expected NextPage 2, got %d", pagination.NextPage)
		}
		if pagination.PrevPage != 0 {
			t.Errorf("expected PrevPage 0, got %d", pagination.PrevPage)
		}
	})

	t.Run("successful pagination - middle page", func(t *testing.T) {
		// Mock count query
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(12))

		// Mock paginated query for page 2 (offset 5)
		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(PostsPerPage, 5).
			WillReturnRows(pgxmock.NewRows([]string{"id", "title", "author", "created", "edited", "body", "description"}))

		_, pagination, err := repo.GetPostsPaginated(2)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if pagination.CurrentPage != 2 {
			t.Errorf("expected CurrentPage 2, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 3 {
			t.Errorf("expected TotalPages 3, got %d", pagination.TotalPages)
		}
		if !pagination.HasNext {
			t.Error("expected HasNext to be true")
		}
		if !pagination.HasPrev {
			t.Error("expected HasPrev to be true")
		}
		if pagination.NextPage != 3 {
			t.Errorf("expected NextPage 3, got %d", pagination.NextPage)
		}
		if pagination.PrevPage != 1 {
			t.Errorf("expected PrevPage 1, got %d", pagination.PrevPage)
		}
	})

	t.Run("count query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnError(pgx.ErrTxClosed)

		posts, pagination, err := repo.GetPostsPaginated(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if posts != nil {
			t.Error("expected nil posts, got non-nil")
		}
		if pagination != (PaginationInfo{}) {
			t.Error("expected empty pagination info")
		}
		if err != nil && !contains(err.Error(), "error getting total posts count") {
			t.Errorf("expected error to contain 'error getting total posts count', got %v", err.Error())
		}
	})

	t.Run("paginated query fails", func(t *testing.T) {
		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM public\.posts`).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(10))

		mock.ExpectQuery(`SELECT \* FROM public\.posts ORDER BY created DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(PostsPerPage, 0).
			WillReturnError(pgx.ErrTxClosed)

		posts, pagination, err := repo.GetPostsPaginated(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if posts != nil {
			t.Error("expected nil posts, got non-nil")
		}
		if pagination != (PaginationInfo{}) {
			t.Error("expected empty pagination info")
		}
		if err != nil && !contains(err.Error(), "error getting paginated posts") {
			t.Errorf("expected error to contain 'error getting paginated posts', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_CreatePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful create", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO public\.posts \(title, description, body, author, created, edited\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\)`).
			WithArgs("New Post", "New description", "new-post.html", "Adam Shkolnik").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := repo.CreatePost("New Post", "New description", "new-post.html", "Adam Shkolnik")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO public\.posts \(title, description, body, author, created, edited\) VALUES \(\$1, \$2, \$3, \$4, NOW\(\), NOW\(\)\)`).
			WithArgs("New Post", "New description", "new-post.html", "Adam Shkolnik").
			WillReturnError(pgx.ErrTxClosed)

		err := repo.CreatePost("New Post", "New description", "new-post.html", "Adam Shkolnik")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != nil && !contains(err.Error(), "error creating post") {
			t.Errorf("expected error to contain 'error creating post', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_UpdatePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectExec(`UPDATE public\.posts SET title = \$2, description = \$3, body = \$4, edited = NOW\(\) WHERE id = \$1`).
			WithArgs(1, "Updated Title", "Updated description", "updated-post.html").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdatePost(1, "Updated Title", "Updated description", "updated-post.html")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mock.ExpectExec(`UPDATE public\.posts SET title = \$2, description = \$3, body = \$4, edited = NOW\(\) WHERE id = \$1`).
			WithArgs(999, "Updated Title", "Updated description", "updated-post.html").
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.UpdatePost(999, "Updated Title", "Updated description", "updated-post.html")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != nil && !contains(err.Error(), "post with id 999 not found") {
			t.Errorf("expected error to contain 'post with id 999 not found', got %v", err.Error())
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(`UPDATE public\.posts SET title = \$2, description = \$3, body = \$4, edited = NOW\(\) WHERE id = \$1`).
			WithArgs(1, "Updated Title", "Updated description", "updated-post.html").
			WillReturnError(pgx.ErrTxClosed)

		err := repo.UpdatePost(1, "Updated Title", "Updated description", "updated-post.html")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != nil && !contains(err.Error(), "error updating post") {
			t.Errorf("expected error to contain 'error updating post', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestConcreteRepository_DeletePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := ConcreteRepository{Pool: mock}

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM public\.posts WHERE id = \$1`).
			WithArgs(1).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.DeletePost(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("post not found", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM public\.posts WHERE id = \$1`).
			WithArgs(999).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.DeletePost(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != nil && !contains(err.Error(), "post with id 999 not found") {
			t.Errorf("expected error to contain 'post with id 999 not found', got %v", err.Error())
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM public\.posts WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(pgx.ErrTxClosed)

		err := repo.DeletePost(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != nil && !contains(err.Error(), "error deleting post") {
			t.Errorf("expected error to contain 'error deleting post', got %v", err.Error())
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
