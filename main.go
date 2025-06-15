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
    "website/internal/database"
    "website/internal/handlers"
    "website/internal/middleware"
    "website/internal/parse"
    "website/internal/posts"
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

    env := handlers.Env{PostsRepository: repo, Templates: templates}

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
    router.HandleFunc("GET /contact", env.ContactHandler)

    router.Handle("GET /static/", http.StripPrefix("/static/", fs))

    if err := server.ListenAndServe(); err != nil {
        return err
    }

    return nil
}
