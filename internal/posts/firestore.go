package posts

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreRepository struct {
	Client     *firestore.Client
	Collection string
}

func NewFirestoreRepository(client *firestore.Client) *FirestoreRepository {
	return &FirestoreRepository{
		Client:     client,
		Collection: "posts",
	}
}

func (repo *FirestoreRepository) GetPost(id int) (*Post, error) {
	ctx := context.Background()
	doc, err := repo.Client.Collection(repo.Collection).Doc(strconv.Itoa(id)).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}

	var post Post
	if err := doc.DataTo(&post); err != nil {
		return nil, fmt.Errorf("error unmarshaling post: %w", err)
	}

	post.ID = id
	return &post, nil
}

func (repo *FirestoreRepository) GetPosts() ([]Post, error) {
	ctx := context.Background()
	iter := repo.Client.Collection(repo.Collection).Documents(ctx)
	defer iter.Stop()

	var posts []Post
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating posts: %w", err)
		}

		var post Post
		if err := doc.DataTo(&post); err != nil {
			return nil, fmt.Errorf("error unmarshaling post: %w", err)
		}

		// Parse ID from document ID
		if id, err := strconv.Atoi(doc.Ref.ID); err == nil {
			post.ID = id
		}

		posts = append(posts, post)
	}

	// Sort by created date descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created.After(posts[j].Created)
	})

	return posts, nil
}

func (repo *FirestoreRepository) GetPostsPaginated(page int) ([]Post, PaginationInfo, error) {
	// Get all posts first to calculate pagination
	allPosts, err := repo.GetPosts()
	if err != nil {
		return nil, PaginationInfo{}, fmt.Errorf("error getting posts for pagination: %w", err)
	}

	totalPosts := len(allPosts)
	paginationInfo := NewPaginationInfo(totalPosts, page)

	// Calculate slice bounds
	offset := paginationInfo.GetOffset()
	end := offset + PostsPerPage
	if end > totalPosts {
		end = totalPosts
	}

	var paginatedPosts []Post
	if offset < totalPosts {
		paginatedPosts = allPosts[offset:end]
	}

	return paginatedPosts, paginationInfo, nil
}

func (repo *FirestoreRepository) GetTotalPostsCount() (int, error) {
	ctx := context.Background()

	// Use Firestore aggregation query for efficient counting
	query := repo.Client.Collection(repo.Collection)
	countQuery := query.NewAggregationQuery().WithCount("total")

	results, err := countQuery.Get(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting post count: %w", err)
	}

	countResult, ok := results["total"]
	if !ok {
		return 0, fmt.Errorf("count result not found")
	}

	countValue, ok := countResult.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected count type: %T", countResult)
	}

	return int(countValue), nil
}

func (repo *FirestoreRepository) DeletePost(id int) error {
	ctx := context.Background()
	docID := strconv.Itoa(id)

	// Check if document exists first
	_, err := repo.Client.Collection(repo.Collection).Doc(docID).Get(ctx)
	if err != nil {
		return fmt.Errorf("post with id %d not found", id)
	}

	_, err = repo.Client.Collection(repo.Collection).Doc(docID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	return nil
}

func (repo *FirestoreRepository) UpdatePost(id int, title, description, body string) error {
	ctx := context.Background()
	docID := strconv.Itoa(id)

	// Check if document exists first
	doc, err := repo.Client.Collection(repo.Collection).Doc(docID).Get(ctx)
	if err != nil {
		return fmt.Errorf("post with id %d not found", id)
	}

	var existingPost Post
	if err := doc.DataTo(&existingPost); err != nil {
		return fmt.Errorf("error reading existing post: %w", err)
	}

	updates := []firestore.Update{
		{Path: "title", Value: title},
		{Path: "description", Value: description},
		{Path: "body", Value: body},
		{Path: "edited", Value: time.Now()},
	}

	_, err = repo.Client.Collection(repo.Collection).Doc(docID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}

	return nil
}

func (repo *FirestoreRepository) CreatePost(title, description, body, author string) error {
	ctx := context.Background()

	// Get next available ID
	nextID, err := repo.getNextID()
	if err != nil {
		return fmt.Errorf("error getting next ID: %w", err)
	}

	now := time.Now()
	post := map[string]interface{}{
		"title":       title,
		"description": description,
		"body":        body,
		"author":      author,
		"created":     now,
		"edited":      now,
	}

	_, err = repo.Client.Collection(repo.Collection).Doc(strconv.Itoa(nextID)).Set(ctx, post)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}

	return nil
}

func (repo *FirestoreRepository) getNextID() (int, error) {
	ctx := context.Background()
	iter := repo.Client.Collection(repo.Collection).Documents(ctx)
	defer iter.Stop()

	maxID := 0
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error iterating documents: %w", err)
		}

		if id, err := strconv.Atoi(doc.Ref.ID); err == nil && id > maxID {
			maxID = id
		}
	}

	return maxID + 1, nil
}
