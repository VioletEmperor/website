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

	"cloud.google.com/go/storage"
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
	} else if conf.StorageMode == "gcs" {
		if conf.GCSBucketName == "" {
			log.Fatal("GCS_BUCKET_NAME is required when using GCS storage mode")
		}

		// Create GCS client
		gcsClient, err := storage.NewClient(ctx)
		if err != nil {
			return err
		}

		contentService = content.NewGCSService(gcsClient, conf.GCSBucketName, conf.GCSPrefix)
		log.Printf("Using GCS content service with bucket: %s, prefix: %s", conf.GCSBucketName, conf.GCSPrefix)
	} else {
		log.Printf("Unknown storage mode '%s', falling back to local filesystem", conf.StorageMode)
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

	// Create separate routers
	publicRouter := http.NewServeMux()
	adminRouter := http.NewServeMux()

	// Middleware stacks
	publicMid := middleware.Stack(middleware.EnableCors, middleware.Logger)
	adminMid := middleware.Stack(middleware.EnableCors, middleware.Logger, middleware.Auth(authClient))

	// Create main router that delegates to sub-routers
	mainRouter := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: mainRouter,
	}

	fs := http.FileServer(http.Dir("static"))

	// Admin routes - relative paths since mounted under /admin/
	adminRouter.HandleFunc("GET /", env.AdminHandler)
	adminRouter.HandleFunc("GET /login", env.AdminLoginPageHandler)
	adminRouter.HandleFunc("POST /logout", env.AdminLogoutHandler)
	adminRouter.HandleFunc("GET /dashboard", env.AdminDashboardHandler)
	adminRouter.HandleFunc("POST /verify", env.AdminVerifyHandler)
	
	// Admin post management routes
	adminRouter.HandleFunc("GET /posts", env.AdminListPostsHandler)
	adminRouter.HandleFunc("GET /posts/{id}", env.AdminGetPostHandler)
	adminRouter.HandleFunc("PUT /posts/{id}", env.AdminUpdatePostHandler)
	adminRouter.HandleFunc("DELETE /posts/{id}", env.AdminDeletePostHandler)
	adminRouter.HandleFunc("POST /posts/upload", env.AdminUploadPostHandler)

	// Public routes - use specific patterns to avoid conflicts
	publicRouter.HandleFunc("GET /{$}", env.RootHandler)
	publicRouter.HandleFunc("GET /about", env.AboutHandler)
	publicRouter.HandleFunc("GET /blog/posts", env.PostsHandler)
	publicRouter.HandleFunc("GET /blog/post/{id}", env.PostHandler)
	publicRouter.HandleFunc("GET /contact", env.ContactHandler)
	publicRouter.HandleFunc("POST /contact", env.MessageHandler)

	// Mount routers with their middleware - strip prefix for admin routes
	mainRouter.Handle("/admin/", http.StripPrefix("/admin", adminMid(adminRouter)))
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fs))
	mainRouter.Handle("/", publicMid(publicRouter))

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
