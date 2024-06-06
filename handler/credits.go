package handler

import (
	"net/http"
	"os"

	"github.com/McFlanky/dreampic-ai/view/credits"
	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, credits.Index())
}

func HandleStripeCheckoutCreate(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:7331/checkout/success/{CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:7331/checkout/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(chi.URLParam(r, "productID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}

func HandleStripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	// sessionID := chi.URLParam(r, "sessionID")
	// stripe.Key = os.Getenv("STRIPE_API_KEY")
	// session, err := session.Get(sessionID, nil)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func HandleStripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}
