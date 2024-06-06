package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/McFlanky/dreampic-ai/db"
	"github.com/McFlanky/dreampic-ai/handler"
	"github.com/McFlanky/dreampic-ai/pkg/sb"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", public())
	// GET
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	// POST
	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Post("/replicate/callback/{userID}/{batchID}", handler.Make(handler.HandleReplicateCallback))

	router.Get("/long-process", handler.Make(handler.HandleLongProcess))

	// With Auth ONLY
	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth)
		auth.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		auth.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))
	})

	// With Auth & Account Setup
	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth, handler.WithAccountSetup)

		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		auth.Get("/generate", handler.Make(handler.HandleGenerateIndex))
		auth.Get("/generate/image/status/{id}", handler.Make(handler.HandleGenerateImageStatus))
		auth.Get("/buy-credits", handler.Make(handler.HandleCreditsIndex))
		auth.Get("/checkout/create/{productID}", handler.Make(handler.HandleStripeCheckoutCreate))
		auth.Get("/checkout/success/{sessionID}", handler.Make(handler.HandleStripeCheckoutSuccess))
		auth.Get("/checkout/cancel", handler.Make(handler.HandleStripeCheckoutCancel))

		auth.Put("/settings/account/profile", handler.Make(handler.HandleSettingsUsernameUpdate))

		auth.Post("/generate", handler.Make(handler.HandleGenerateCreate))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("Starting server on ", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
