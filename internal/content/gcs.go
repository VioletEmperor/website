package content

import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
)

// GCSService implements ContentService for Google Cloud Storage
type GCSService struct {
	client     *storage.Client
	bucketName string
	prefix     string // Optional prefix for post files (e.g., "posts/")
}

// NewGCSService creates a new GCS-based content service
func NewGCSService(client *storage.Client, bucketName, prefix string) *GCSService {
	return &GCSService{
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}
}

// GetContent retrieves HTML content from Google Cloud Storage
func (gcs *GCSService) GetContent(filename string) (string, error) {
	// Security check: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		return "", fmt.Errorf("invalid filename: %s", filename)
	}
	
	// Construct the object path
	objectPath := filename
	if gcs.prefix != "" {
		objectPath = gcs.prefix + filename
	}
	
	// Get the object from GCS
	ctx := context.Background()
	bucket := gcs.client.Bucket(gcs.bucketName)
	obj := bucket.Object(objectPath)
	
	// Open a reader for the object
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return "", fmt.Errorf("post content not found: %s", filename)
		}
		return "", fmt.Errorf("failed to open GCS object: %w", err)
	}
	defer reader.Close()
	
	// Read the content
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read GCS object content: %w", err)
	}
	
	return string(content), nil
}