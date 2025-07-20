package handlers

import (
	"html/template"
	"website/internal/config"
	"website/internal/posts"

	"firebase.google.com/go/v4/auth"
)

type Env struct {
	PostsRepository posts.Repository
	Templates       map[string]*template.Template
	EmailKey        string
	FirebaseAuth    *auth.Client
	Config          config.Config
}
