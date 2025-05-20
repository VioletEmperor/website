package main

import (
    "log"
    "net/http"
    "website/posts_repository"
)

func main() {
    log.Println("starting server...")

    config, err := getConfig()

    if err != nil {
        log.Fatal(err)
    }

    templates := parse()

    pool, err := connect(config.URL)

    if err != nil {
        log.Fatal(err)
    }

    defer pool.Close()

    repo := posts_repository.New(pool)

    router := http.NewServeMux()

    server := &http.Server{
        Addr:    ":" + config.Port,
        Handler: router,
    }

    fs := http.FileServer(http.Dir("static"))

    router.HandleFunc("/", rootHandler)
    router.HandleFunc("/about", aboutHandler(templates))
    router.HandleFunc("/posts", postsHandler(repo, templates))
    router.HandleFunc("/contact", contactHandler(templates))
    router.Handle("/static/", http.StripPrefix("/static/", fs))

    log.Fatal(server.ListenAndServe())
}
