package content

import (
	"fmt"
	"os"
	"path/filepath"
)

// FilesystemService implements ContentService for local filesystem storage
type FilesystemService struct {
	postsDirectory string
}

// NewFilesystemService creates a new filesystem-based content service
func NewFilesystemService(postsDirectory string) *FilesystemService {
	return &FilesystemService{
		postsDirectory: postsDirectory,
	}
}

// GetContent retrieves HTML content from the local filesystem
func (fs *FilesystemService) GetContent(filename string) (string, error) {
	// Construct the full file path
	filePath := filepath.Join(fs.postsDirectory, filename)
	
	// Security check: ensure the file is within the posts directory
	absPostsDir, err := filepath.Abs(fs.postsDirectory)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for posts directory: %w", err)
	}
	
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for file: %w", err)
	}
	
	// Check if the file is within the posts directory (prevent directory traversal)
	if !filepath.HasPrefix(absFilePath, absPostsDir) {
		return "", fmt.Errorf("file path outside posts directory: %s", filename)
	}
	
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("post content not found: %s", filename)
		}
		return "", fmt.Errorf("failed to read post content: %w", err)
	}
	
	return string(content), nil
}