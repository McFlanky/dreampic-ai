package handler

import (
	"net/http"

	"github.com/McFlanky/dreampic-ai/view/credits"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, credits.Index())
}
