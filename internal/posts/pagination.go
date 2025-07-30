package posts

const PostsPerPage = 5

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	CurrentPage int
	TotalPages  int
	TotalPosts  int
	HasNext     bool
	HasPrev     bool
	NextPage    int
	PrevPage    int
}

// NewPaginationInfo creates pagination info from total posts and current page
func NewPaginationInfo(totalPosts, currentPage int) PaginationInfo {
	// Calculate total pages (round up)
	totalPages := (totalPosts + PostsPerPage - 1) / PostsPerPage
	
	// Ensure currentPage is within valid bounds
	if currentPage < 1 {
		currentPage = 1
	}
	if totalPages == 0 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}
	
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1
	
	nextPage := currentPage + 1
	if !hasNext {
		nextPage = 0 // Indicate no next page
	}
	
	prevPage := currentPage - 1
	if !hasPrev {
		prevPage = 0 // Indicate no previous page
	}
	
	return PaginationInfo{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalPosts:  totalPosts,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

// GetOffset calculates the SQL OFFSET for the current page
func (p PaginationInfo) GetOffset() int {
	return (p.CurrentPage - 1) * PostsPerPage
}