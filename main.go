package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"website/internal/config"
	"website/internal/content"
	"website/internal/database"
	"website/internal/handlers"
	"website/internal/middleware"
	"website/internal/parse"
	"website/internal/posts"

	firebase "firebase.google.com/go/v4"
)

func main() {
	log.Println("starting server...")

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := run(ctx, cancel); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit

	log.Println("shutting down server...")
}

func run(ctx context.Context, cancel context.CancelFunc) error {
	defer cancel()

	conf, err := config.GetConfig()

	if err != nil {
		return err
	}

	templates := parse.Parse()

	pool, err := database.Connect(ctx, conf.URL)

	if err != nil {
		return err
	}

	repo := posts.New(pool)

	// Initialize content service based on storage mode
	var contentService content.ContentService
	if conf.StorageMode == "local" {
		contentService = content.NewFilesystemService(conf.PostsDirectory)
		log.Printf("Using local filesystem content service with directory: %s", conf.PostsDirectory)
	} else {
		// TODO: Implement GCS service when needed
		log.Printf("GCS storage mode not yet implemented, falling back to local filesystem")
		contentService = content.NewFilesystemService(conf.PostsDirectory)
	}

	// Initialize Firebase Auth
	firebaseConf := &firebase.Config{
		ProjectID: conf.ProjectID,
	}

	app, err := firebase.NewApp(ctx, firebaseConf)
	if err != nil {
		return err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return err
	}

	env := handlers.Env{
		PostsRepository: repo,
		ContentService:  contentService,
		Templates:       templates,
		EmailKey:        conf.EmailKey,
		FirebaseAuth:    authClient,
		Config:          conf,
	}

	router := http.NewServeMux()

	mid := middleware.Stack(middleware.EnableCors, middleware.Logger)

	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mid(router),
	}

	fs := http.FileServer(http.Dir("static"))

	router.HandleFunc("GET /", env.RootHandler)
	router.HandleFunc("GET /about", env.AboutHandler)
	router.HandleFunc("GET /blog/posts", env.PostsHandler)
	router.HandleFunc("GET /blog/post/{id}", env.PostHandler)
	router.HandleFunc("GET /contact", env.ContactHandler)
	router.HandleFunc("POST /contact", env.MessageHandler)

	// Admin routes
	router.HandleFunc("GET /admin", env.AdminHandler)
	router.HandleFunc("GET /admin/login", env.AdminLoginPageHandler)
	router.HandleFunc("POST /admin/logout", env.AdminLogoutHandler)
	router.HandleFunc("GET /admin/dashboard", env.AdminDashboardHandler)
	router.HandleFunc("POST /admin/verify", env.AdminVerifyHandler)

	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
