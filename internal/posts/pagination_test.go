package posts

import (
	"testing"
)

func TestNewPaginationInfo(t *testing.T) {
	t.Run("first page with multiple pages", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 1)

		if pagination.CurrentPage != 1 {
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 3 { // 12 posts / 5 per page = 3 pages
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

	t.Run("middle page", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 2)

		if pagination.CurrentPage != 2 {
			t.Errorf("expected CurrentPage 2, got %d", pagination.CurrentPage)
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

	t.Run("last page", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 3)

		if pagination.CurrentPage != 3 {
			t.Errorf("expected CurrentPage 3, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 3 {
			t.Errorf("expected TotalPages 3, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 12 {
			t.Errorf("expected TotalPosts 12, got %d", pagination.TotalPosts)
		}
		if pagination.HasNext {
			t.Error("expected HasNext to be false")
		}
		if !pagination.HasPrev {
			t.Error("expected HasPrev to be true")
		}
		if pagination.NextPage != 0 {
			t.Errorf("expected NextPage 0, got %d", pagination.NextPage)
		}
		if pagination.PrevPage != 2 {
			t.Errorf("expected PrevPage 2, got %d", pagination.PrevPage)
		}
	})

	t.Run("single page", func(t *testing.T) {
		pagination := NewPaginationInfo(3, 1)

		if pagination.CurrentPage != 1 {
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 1 {
			t.Errorf("expected TotalPages 1, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 3 {
			t.Errorf("expected TotalPosts 3, got %d", pagination.TotalPosts)
		}
		if pagination.HasNext {
			t.Error("expected HasNext to be false")
		}
		if pagination.HasPrev {
			t.Error("expected HasPrev to be false")
		}
		if pagination.NextPage != 0 {
			t.Errorf("expected NextPage 0, got %d", pagination.NextPage)
		}
		if pagination.PrevPage != 0 {
			t.Errorf("expected PrevPage 0, got %d", pagination.PrevPage)
		}
	})

	t.Run("no posts", func(t *testing.T) {
		pagination := NewPaginationInfo(0, 1)

		if pagination.CurrentPage != 1 {
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 0 {
			t.Errorf("expected TotalPages 0, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 0 {
			t.Errorf("expected TotalPosts 0, got %d", pagination.TotalPosts)
		}
		if pagination.HasNext {
			t.Error("expected HasNext to be false")
		}
		if pagination.HasPrev {
			t.Error("expected HasPrev to be false")
		}
		if pagination.NextPage != 0 {
			t.Errorf("expected NextPage 0, got %d", pagination.NextPage)
		}
		if pagination.PrevPage != 0 {
			t.Errorf("expected PrevPage 0, got %d", pagination.PrevPage)
		}
	})

	t.Run("exact multiple of posts per page", func(t *testing.T) {
		pagination := NewPaginationInfo(10, 1) // 10 posts / 5 per page = exactly 2 pages

		if pagination.CurrentPage != 1 {
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 2 {
			t.Errorf("expected TotalPages 2, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 10 {
			t.Errorf("expected TotalPosts 10, got %d", pagination.TotalPosts)
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

	t.Run("page number less than 1", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 0)

		if pagination.CurrentPage != 1 { // Should be corrected to 1
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
	})

	t.Run("page number greater than total pages", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 5)

		if pagination.CurrentPage != 3 { // Should be corrected to last page
			t.Errorf("expected CurrentPage 3, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 3 {
			t.Errorf("expected TotalPages 3, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 12 {
			t.Errorf("expected TotalPosts 12, got %d", pagination.TotalPosts)
		}
		if pagination.HasNext {
			t.Error("expected HasNext to be false")
		}
		if !pagination.HasPrev {
			t.Error("expected HasPrev to be true")
		}
	})

	t.Run("page number greater than total pages with no posts", func(t *testing.T) {
		pagination := NewPaginationInfo(0, 5)

		if pagination.CurrentPage != 1 { // Should default to 1 when no pages exist
			t.Errorf("expected CurrentPage 1, got %d", pagination.CurrentPage)
		}
		if pagination.TotalPages != 0 {
			t.Errorf("expected TotalPages 0, got %d", pagination.TotalPages)
		}
		if pagination.TotalPosts != 0 {
			t.Errorf("expected TotalPosts 0, got %d", pagination.TotalPosts)
		}
		if pagination.HasNext {
			t.Error("expected HasNext to be false")
		}
		if pagination.HasPrev {
			t.Error("expected HasPrev to be false")
		}
	})
}

func TestPaginationInfo_GetOffset(t *testing.T) {
	t.Run("first page offset", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 1)
		if pagination.GetOffset() != 0 {
			t.Errorf("expected offset 0, got %d", pagination.GetOffset())
		}
	})

	t.Run("second page offset", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 2)
		if pagination.GetOffset() != 5 { // (2-1) * 5 = 5
			t.Errorf("expected offset 5, got %d", pagination.GetOffset())
		}
	})

	t.Run("third page offset", func(t *testing.T) {
		pagination := NewPaginationInfo(12, 3)
		if pagination.GetOffset() != 10 { // (3-1) * 5 = 10
			t.Errorf("expected offset 10, got %d", pagination.GetOffset())
		}
	})
}

func TestPostsPerPage(t *testing.T) {
	// Ensure PostsPerPage constant is set correctly
	if PostsPerPage != 5 {
		t.Errorf("expected PostsPerPage 5, got %d", PostsPerPage)
	}
}