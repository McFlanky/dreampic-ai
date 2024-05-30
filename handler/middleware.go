package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/McFlanky/dreampic-ai/pkg/sb"
	"github.com/McFlanky/dreampic-ai/types"
)

func WithoutAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)
		if !user.IsLoggedIn {
			path := r.URL.Path
			http.Redirect(w, r, "/login?to="+path, http.StatusSeeOther)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := sb.Client.Auth.User(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			Email:      resp.Email,
			IsLoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(fn)
}
