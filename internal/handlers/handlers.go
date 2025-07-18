package handlers

import (
	"fmt"
	"github.com/resend/resend-go/v2"
	"log"
	"net/http"
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

func (env Env) MessageHandler(w http.ResponseWriter, r *http.Request) {
	type Form struct {
		Name    string
		Email   string
		Subject string
		Message string
	}

	w.Header().Set("Content-Type", "text/html; text/css; application/javascript; charset=utf-8")

	message := Form{}

	message.Name = r.FormValue("name")
	message.Email = r.FormValue("email")
	message.Subject = r.FormValue("subject")
	message.Message = r.FormValue("message")

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
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to send email:", err)
		return
	}

	log.Println("sent: ", sent.Id)

	if err := env.Templates["submit.html"].ExecuteTemplate(w, "submit.html", message); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed to execute template:", err)
		return
	}
}
