package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	input := replicate.PredictionInput{
		"prompt": "sexy woman bending over ass spread open for camera",
	}

	webhook := replicate.Webhook{
		URL:    "https://webhook.site/b5a903c3-e4f6-4391-aff3-6ac7ee06071b",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, err := r8.CreatePrediction(ctx, "ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4", input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ", output)

	// // Create a prediction
	// prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	// if err != nil {
	// 	// handle error
	// }
}
