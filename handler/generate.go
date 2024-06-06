package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/McFlanky/dreampic-ai/db"
	"github.com/McFlanky/dreampic-ai/pkg/kit/validate"
	"github.com/McFlanky/dreampic-ai/types"
	"github.com/McFlanky/dreampic-ai/view/generate"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
)

const creditsPerImage = 2

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
	}
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	amount, _ := strconv.Atoi(r.FormValue("amount"))
	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}
	var errors generate.FormErrors

	if amount <= 0 || amount > 8 {
		errors.Amount = "Please enter a valid amount"
		return render(r, w, generate.Form(params, errors))
	}

	ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Min(10), validate.Max(100)),
	}).Validate(&errors)
	if !ok {
		return render(r, w, generate.Form(params, errors))
	}

	creditsNeeded := params.Amount * creditsPerImage
	if user.Account.Credits < creditsNeeded {
		errors.CreditsNeeded = creditsNeeded
		errors.UserCredits = user.Account.Credits
		errors.Credits = true
		return render(r, w, generate.Form(params, errors))
	}

	batchID := uuid.New()
	genParams := GenerateImageParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.ID,
		BatchID: batchID,
	}

	if err := generateImages(r.Context(), genParams); err != nil {
		return err
	}

	err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i := 0; i < params.Amount; i++ {
			img := types.Image{
				UserID:  user.ID,
				Status:  types.ImageStatusPending,
				BatchID: batchID,
			}
			if err := db.CreateImage(tx, &img); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/generate")
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := db.GetImageByID(id)
	if err != nil {
		return err
	}
	slog.Info("checking image status", "id", id)
	return render(r, w, generate.GalleryImage(image))
}

type GenerateImageParams struct {
	Prompt  string
	Amount  int
	BatchID uuid.UUID
	UserID  uuid.UUID
}

func generateImages(ctx context.Context, params GenerateImageParams) error {
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	input := replicate.PredictionInput{
		"prompt":      params.Prompt,
		"num_outputs": params.Amount,
	}
	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("https://webhook.site/b5a903c3-e4f6-4391-aff3-6ac7ee06071b/%s/%s", params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}
	_, err = r8.CreatePrediction(ctx, "ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4", input, &webhook, false)
	return err
}
