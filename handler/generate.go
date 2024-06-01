package handler

import (
	"net/http"

	"github.com/McFlanky/dreampic-ai/view/generate"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, generate.Index())
}
