package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	const PORT = "8080"
	const URL = "postgres://docker:password@database:5432/database?sslmode=disable"

	if err := os.Setenv("PORT", PORT); err != nil {
		t.Errorf("tried to set environment variable PORT for test but got %v", err)
	}

	if err := os.Setenv("URL", URL); err != nil {
		t.Errorf("tried to set environment variable URL for test but got %v", err)
	}

	expected := Config{Port: PORT, URL: URL}

	res, err := GetConfig()

	if err != nil {
		t.Errorf("tried to get config but got %v", err)
	}

	if res != expected {
		t.Errorf("expected %v but got %v", expected, res)
	}
}
