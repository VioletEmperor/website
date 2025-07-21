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
}

func GetConfig() (Config, error) {
	config := Config{}

	config.Port = os.Getenv("PORT")

	if config.Port == "" {
		return config, errors.New("missing environment variable PORT")
	}

	host := os.Getenv("DB_HOST")

	if host == "" {
		return config, errors.New("missing environment variable DB_HOST")
	}

	user := os.Getenv("DB_USER")

	if user == "" {
		return config, errors.New("missing environment variable DB_USER")
	}

	password := os.Getenv("DB_PASSWORD")

	if password == "" {
		return config, errors.New("missing environment variable DB_PASSWORD")
	}

	name := os.Getenv("DB_NAME")

	if name == "" {
		return config, errors.New("missing environment variable DB_NAME")
	}

	config.URL = fmt.Sprintf("user=%s password=%s database=%s host=%s", user, password, name, host)

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

	// Content storage configuration
	config.PostsDirectory = os.Getenv("POSTS_DIRECTORY")
	if config.PostsDirectory == "" {
		config.PostsDirectory = "posts" // Default to posts/ directory
	}

	config.StorageMode = os.Getenv("STORAGE_MODE")
	if config.StorageMode == "" {
		config.StorageMode = "local" // Default to local filesystem
	}

	log.Println(config.URL)

	return config, nil
}
