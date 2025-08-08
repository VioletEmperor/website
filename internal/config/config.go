package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port              string
	URL               string
	EmailKey          string
	ProjectID         string
	FirebaseWebAPIKey string
	PostsDirectory    string
	StorageMode       string // "local" or "gcs"
	GCSBucketName     string
	GCSPrefix         string
}

func GetConfig() (Config, error) {
	config := Config{}

	config.Port = os.Getenv("PORT")

	if config.Port == "" {
		return config, errors.New("missing environment variable PORT")
	}

	// Content storage configuration
	config.StorageMode = os.Getenv("STORAGE_MODE")
	if config.StorageMode == "" {
		config.StorageMode = "local" // Default to local filesystem
	}

	// Database configuration (only required for local/PostgreSQL mode)
	if config.StorageMode == "local" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			return config, errors.New("missing environment variable DB_HOST (required for local storage mode)")
		}

		user := os.Getenv("DB_USER")
		if user == "" {
			return config, errors.New("missing environment variable DB_USER (required for local storage mode)")
		}

		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			return config, errors.New("missing environment variable DB_PASSWORD (required for local storage mode)")
		}

		name := os.Getenv("DB_NAME")
		if name == "" {
			return config, errors.New("missing environment variable DB_NAME (required for local storage mode)")
		}

		config.URL = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, password, name, host)
	}

	config.EmailKey = os.Getenv("EMAIL_KEY")

	if config.EmailKey == "" {
		return config, errors.New("missing environment variable EMAIL_KEY")
	}

	config.ProjectID = os.Getenv("PROJECT_ID")

	if config.ProjectID == "" {
		return config, errors.New("missing environment variable PROJECT_ID")
	}

	config.FirebaseWebAPIKey = os.Getenv("FIREBASE_WEB_API_KEY")

	if config.FirebaseWebAPIKey == "" {
		return config, errors.New("missing environment variable FIREBASE_WEB_API_KEY")
	}

	// Posts directory configuration
	config.PostsDirectory = os.Getenv("POSTS_DIRECTORY")
	if config.PostsDirectory == "" {
		config.PostsDirectory = "posts" // Default to posts/ directory
	}

	// GCS configuration (only required if using GCS storage mode)
	if config.StorageMode == "gcs" {
		config.GCSBucketName = os.Getenv("GCS_BUCKET_NAME")
		config.GCSPrefix = os.Getenv("GCS_PREFIX") // Optional prefix, e.g., "posts/"
	}

	log.Println(config.URL)

	return config, nil
}
