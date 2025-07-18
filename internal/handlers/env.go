package handlers

import (
	"html/template"
	"website/internal/posts"
)

type Env struct {
	PostsRepository posts.Repository
	Templates       map[string]*template.Template
	EmailKey        string
}
