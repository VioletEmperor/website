package content

// ContentService defines the interface for retrieving and storing blog post content
type ContentService interface {
	// GetContent retrieves the HTML content for a blog post by filename
	// Returns the raw HTML content or an error if the file cannot be found/read
	GetContent(filename string) (string, error)
	
	// SaveContent saves HTML content to storage with the given filename
	// Returns an error if the content cannot be saved
	SaveContent(filename, content string) error
}