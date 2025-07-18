package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port     string
	URL      string
	EmailKey string
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

	log.Println(config.URL)

	return config, nil
}
