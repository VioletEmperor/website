package main

import (
    "html/template"
    "log"
    "net/http"
    "website/posts_repository"
)

type PostsPage struct {
    Posts  []posts_repository.Post
    Active string
}

func postsHandler(repo posts_repository.PostsRepository, templates map[string]*template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL)

        w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

        posts, err := repo.GetPosts()

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Println("failed to fetch posts:", err)
            return
        }

        err = templates["posts.html"].ExecuteTemplate(w, "posts.html", PostsPage{posts, "posts"})

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Println("failed to execute template:", err)
            return
        }
    }
}

type AboutPage struct {
    Active string
}

func aboutHandler(templates map[string]*template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL)

        w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

        err := templates["about.html"].ExecuteTemplate(w, "about.html", AboutPage{"about"})

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Println("failed to execute template:", err)
            return
        }
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/about", http.StatusFound)
}

type ContactPage struct {
    Active string
}

type Form struct {
    Name    string
    Email   string
    Message string
}

func contactHandler(templates map[string]*template.Template) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL, r.Method)

        w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

        if r.Method == "GET" {
            err := templates["contact.html"].ExecuteTemplate(w, "contact.html", ContactPage{"contact"})

            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Println("failed to execute template:", err)
                return
            }
        } else if r.Method == "POST" {
            message := Form{}

            message.Name = r.FormValue("name")
            message.Email = r.FormValue("email")
            message.Message = r.FormValue("message")

            err := templates["submit.html"].ExecuteTemplate(w, "submit.html", message)

            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                log.Println("failed to execute template:", err)
                return
            }
        } else {

        }
    }
}
