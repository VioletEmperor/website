package handlers

import (
	"context"
	"fmt"
	"github.com/resend/resend-go/v2"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"website/internal/posts"
)

func (env Env) PostsHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Posts  []posts.Post
		Active string
	}

	log.Println(r.Header.Get("Accept"))

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	list, err := env.PostsRepository.GetPosts()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to fetch posts:", err)
		return
	}

	err = env.Templates["posts.html"].ExecuteTemplate(w, "posts.html", Data{list, "posts"})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) PostHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Post    posts.Post
		Content template.HTML
		Active  string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	// Extract post ID from URL path
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := env.PostsRepository.GetPost(id)
	if err != nil {
		log.Printf("failed to fetch post %d: %v", id, err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Load HTML content for the post
	htmlContent, err := env.ContentService.GetContent(post.Body)
	if err != nil {
		log.Printf("failed to load content for post %d (file: %s): %v", id, post.Body, err)
		http.Error(w, "Post content not available", http.StatusNotFound)
		return
	}

	data := Data{
		Post:    *post,
		Content: template.HTML(htmlContent),
		Active:  "posts",
	}

	err = env.Templates["post.html"].ExecuteTemplate(w, "post.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AboutHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Active string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	err := env.Templates["about.html"].ExecuteTemplate(w, "about.html", Data{"about"})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/about", http.StatusFound)
}

func (env Env) ContactHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Active string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	err := env.Templates["contact.html"].ExecuteTemplate(w, "contact.html", Data{"contact"})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to login page first - users need to authenticate before accessing dashboard
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func (env Env) AdminLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Active         string
		FirebaseAPIKey string
		ProjectID      string
		Error          string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	data := Data{
		Active:         "admin",
		FirebaseAPIKey: env.Config.FirebaseWebAPIKey,
		ProjectID:      env.Config.ProjectID,
		Error:          r.URL.Query().Get("error"),
	}

	err := env.Templates["admin-login.html"].ExecuteTemplate(w, "admin-login.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Client-side authentication handles logout
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func (env Env) AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Active         string
		FirebaseAPIKey string
		ProjectID      string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	data := Data{
		Active:         "admin",
		FirebaseAPIKey: env.Config.FirebaseWebAPIKey,
		ProjectID:      env.Config.ProjectID,
	}

	err := env.Templates["admin-dashboard.html"].ExecuteTemplate(w, "admin-dashboard.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AdminVerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>" format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return
	}

	idToken := parts[1]

	// Verify the ID token with Firebase
	_, err := env.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("firebase token verification failed: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("OK"))
}

func (env Env) MessageHandler(w http.ResponseWriter, r *http.Request) {
	type Form struct {
		Name    string
		Email   string
		Subject string
		Message string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	message := Form{}

	message.Name = strings.TrimSpace(r.FormValue("name"))
	message.Email = strings.TrimSpace(r.FormValue("email"))
	message.Subject = strings.TrimSpace(r.FormValue("subject"))
	message.Message = strings.TrimSpace(r.FormValue("message"))

	// Validate required fields
	if message.Name == "" || message.Email == "" || message.Subject == "" || message.Message == "" {
		log.Println("missing required form fields")
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(message.Email); err != nil {
		log.Println("invalid email format:", message.Email)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	client := resend.NewClient(env.EmailKey)

	params := &resend.SendEmailRequest{
		From:        "contact@adamshkolnik.com",
		To:          []string{message.Email},
		Subject:     fmt.Sprintf("Message From %s Has Been Received Successfully!", message.Name),
		Bcc:         nil,
		Cc:          []string{"adam.shkolnik@outlook.com"},
		ReplyTo:     "",
		Html:        fmt.Sprintf("<html><body><strong>%s</strong>\n<p>%s</p></body></html>", message.Subject, message.Message),
		Text:        "",
		Tags:        nil,
		Attachments: nil,
		Headers:     nil,
		ScheduledAt: "",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		// Check for specific error types
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "authentication") || strings.Contains(errorMsg, "unauthorized") {
			log.Println("email service authentication failed:", err)
			http.Error(w, "Email service temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
		if strings.Contains(errorMsg, "rate limit") || strings.Contains(errorMsg, "quota") {
			log.Println("email service rate limited:", err)
			http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
			return
		}
		// Generic error
		log.Println("failed to send email:", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	log.Println("sent: ", sent.Id)

	if err := env.Templates["submit.html"].ExecuteTemplate(w, "submit.html", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}
