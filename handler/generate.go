package handler

import (
	"log/slog"
	"net/http"

	"github.com/McFlanky/dreampic-ai/types"
	"github.com/McFlanky/dreampic-ai/view/generate"
	"github.com/go-chi/chi/v5"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	// images := make([]types.Image, 20)
	data := generate.ViewData{
		Images: []types.Image{},
	}
	// images[0].Status = types.ImageStatusPending
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, generate.GalleryImage(types.Image{Status: types.ImageStatusPending}))
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	// fetch from DB
	image := types.Image{
		Status: types.ImageStatusPending,
	}
	slog.Info("checking image status", "id", id)
	return render(r, w, generate.GalleryImage(image))
}
