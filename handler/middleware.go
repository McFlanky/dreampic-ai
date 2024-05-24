package handler

import (
	"fmt"
	"net/http"
	"strings"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("from the with user middleware")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
