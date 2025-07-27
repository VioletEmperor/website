package handlers

import (
	"firebase.google.com/go/v4/auth"
	"html/template"
	"website/internal/config"
	"website/internal/content"
	"website/internal/posts"
)

type Env struct {
	PostsRepository posts.Repository
	ContentService  content.ContentService
	Templates       map[string]*template.Template
	EmailKey        string
	FirebaseAuth    *auth.Client
	Config          config.Config
}
