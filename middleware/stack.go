package middleware

import "net/http"

func Stack(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for _, m := range middlewares {
        next = m(next)
    }

    return next
}
