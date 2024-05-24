package handler

import (
	"net/http"

	"github.com/McFlanky/dreampic-ai/view/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}
