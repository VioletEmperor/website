package main

import (
    "context"
    "html/template"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "website/middleware"
    "website/posts_repository"
)

func main() {
    log.Println("starting server...")

    ctx, cancel := context.WithCancel(context.Background())

    quit := make(chan os.Signal, 1)

    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := run(ctx, cancel); err != nil {
            log.Fatal(err)
        }
    }()

    <-quit

    log.Println("shutting down server...")
}

type Env struct {
    PostsRepository posts_repository.PostsRepository
    Templates       map[string]*template.Template
}

func run(ctx context.Context, cancel context.CancelFunc) error {
    defer cancel()

    config, err := getConfig()

    if err != nil {
        return err
    }

    templates := parse()

    pool, err := connect(ctx, config.URL)

    if err != nil {
        return err
    }

    repo := posts_repository.New(pool)

    env := Env{repo, templates}

    router := http.NewServeMux()

    server := &http.Server{
        Addr: ":" + config.Port,
        Handler: middleware.Stack(
            router,
            middleware.EnforceJSON,
            middleware.EnableCors,
            middleware.Logger),
    }

    fs := http.FileServer(http.Dir("static"))

    router.HandleFunc("/", env.rootHandler)
    router.HandleFunc("/about", env.aboutHandler)
    router.HandleFunc("/blog/posts", env.postsHandler)
    router.HandleFunc("/contact", env.contactHandler)
    router.Handle("/static/", http.StripPrefix("/static/", fs))

    addRoutes(router)

    if err := server.ListenAndServe(); err != nil {
        return err
    }

    return nil
}
