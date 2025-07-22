package middleware

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"log"
	"net/http"
)

func Auth(firebaseAuth *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for login page, verify endpoint, and root admin (which redirects to login)
			// Note: paths are relative since /admin prefix is stripped
			if r.URL.Path == "/" || r.URL.Path == "/login" || r.URL.Path == "/verify" {
				next.ServeHTTP(w, r)
				return
			}

			// Get token from cookie
			cookie, err := r.Cookie("adminToken")
			if err != nil {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			idToken := cookie.Value

			// Verify the ID token with Firebase
			_, err = firebaseAuth.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				log.Printf("firebase token verification failed: %v", err)
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}

			// Token is valid, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}
