package handlers

import (
	"context"
	"encoding/json"
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
		Posts      []posts.Post
		Pagination posts.PaginationInfo
		Active     string
	}

	log.Println(r.Header.Get("Accept"))

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	// Parse page parameter from query string
	pageStr := r.URL.Query().Get("page")
	page := 1 // Default to page 1
	
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage < 1 {
			// Invalid page parameter, redirect to page 1
			http.Redirect(w, r, "/blog/posts?page=1", http.StatusSeeOther)
			return
		}
		page = parsedPage
	}

	// Get paginated posts
	list, paginationInfo, err := env.PostsRepository.GetPostsPaginated(page)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to fetch paginated posts:", err)
		return
	}

	// If page is out of bounds and we have posts, redirect to last page
	if paginationInfo.TotalPages > 0 && page > paginationInfo.TotalPages {
		http.Redirect(w, r, fmt.Sprintf("/blog/posts?page=%d", paginationInfo.TotalPages), http.StatusSeeOther)
		return
	}

	err = env.Templates["posts.html"].ExecuteTemplate(w, "posts.html", Data{list, paginationInfo, "posts"})

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
		Posts          []posts.Post
	}

	// Get posts for dashboard
	postsList, err := env.PostsRepository.GetPosts()
	if err != nil {
		log.Printf("failed to fetch posts for dashboard: %v", err)
		postsList = []posts.Post{} // Empty slice if error
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	data := Data{
		Active:         "admin",
		FirebaseAPIKey: env.Config.FirebaseWebAPIKey,
		ProjectID:      env.Config.ProjectID,
		Posts:          postsList,
	}

	err = env.Templates["admin-dashboard.html"].ExecuteTemplate(w, "admin-dashboard.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AdminVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if !env.verifyAdminAuth(w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)
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

	if err := env.Templates["partials/submit.html"].ExecuteTemplate(w, "partials/submit.html", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}

func (env Env) AdminDeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication
	if !env.verifyAdminAuth(w, r) {
		return
	}

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

	// Delete the post
	err = env.PostsRepository.DeletePost(id)
	if err != nil {
		log.Printf("failed to delete post %d: %v", id, err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env Env) AdminUpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication
	if !env.verifyAdminAuth(w, r) {
		return
	}

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

	// Parse JSON body
	var updateData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if updateData.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// If body is empty, get the original body to keep it unchanged
	if updateData.Body == "" {
		post, err := env.PostsRepository.GetPost(id)
		if err != nil {
			http.Error(w, "Failed to get original post", http.StatusInternalServerError)
			return
		}
		updateData.Body = post.Body
	}

	// Update the post
	err = env.PostsRepository.UpdatePost(id, updateData.Title, updateData.Description, updateData.Body)
	if err != nil {
		log.Printf("failed to update post %d: %v", id, err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (env Env) AdminGetPostHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication
	if !env.verifyAdminAuth(w, r) {
		return
	}

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

	// Get the post
	post, err := env.PostsRepository.GetPost(id)
	if err != nil {
		log.Printf("failed to get post %d: %v", id, err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (env Env) AdminListPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication
	if !env.verifyAdminAuth(w, r) {
		return
	}

	// Get all posts
	posts, err := env.PostsRepository.GetPosts()
	if err != nil {
		log.Printf("failed to get posts: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (env Env) AdminUploadPostHandler(w http.ResponseWriter, r *http.Request) {
	// Verify authentication
	if !env.verifyAdminAuth(w, r) {
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := strings.TrimSpace(r.FormValue("excerpt"))
	editMode := r.FormValue("editMode")
	postIdStr := r.FormValue("postId")

	// Validate required fields
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("htmlFile")
	if err != nil {
		http.Error(w, "HTML file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasSuffix(header.Filename, ".html") {
		http.Error(w, "Only HTML files are allowed", http.StatusBadRequest)
		return
	}

	// Read file content
	content := make([]byte, header.Size)
	_, err = file.Read(content)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Store the file using the content service
	bodyFilename := header.Filename
	err = env.ContentService.SaveContent(bodyFilename, string(content))
	if err != nil {
		log.Printf("failed to save content file %s: %v", bodyFilename, err)
		http.Error(w, "Failed to save file content", http.StatusInternalServerError)
		return
	}

	if editMode == "true" && postIdStr != "" {
		// Handle edit mode
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Update existing post
		err = env.PostsRepository.UpdatePost(postId, title, excerpt, bodyFilename)
		if err != nil {
			log.Printf("failed to update post %d: %v", postId, err)
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// Redirect back to dashboard
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	} else {
		// Handle new post creation
		// For now, use a default author - you could get this from the authenticated user
		author := "Admin"
		err = env.PostsRepository.CreatePost(title, excerpt, bodyFilename, author)
		if err != nil {
			log.Printf("failed to create new post: %v", err)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// Redirect back to dashboard
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

// Helper function to verify admin authentication
func (env Env) verifyAdminAuth(w http.ResponseWriter, r *http.Request) bool {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return false
	}

	// Extract token from "Bearer <token>" format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
		return false
	}

	idToken := parts[1]

	// Verify the ID token with Firebase
	_, err := env.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		log.Printf("firebase token verification failed: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return false
	}

	return true
}
