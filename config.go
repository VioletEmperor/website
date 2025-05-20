package main

import (
    "errors"
    "os"
)

type Config struct {
    Port string
    URL  string
}

func getConfig() (Config, error) {
    config := Config{}

    config.Port = os.Getenv("PORT")

    if config.Port == "" {
        return config, errors.New("missing environment variable PORT")
    }

    config.URL = os.Getenv("URL")

    if config.URL == "" {
        return config, errors.New("missing environment variable URL")
    }

    return config, nil
}
