package handler

import (
	"log/slog"
	"net/http"

	"github.com/McFlanky/dreampic-ai/pkg/sb"
	"github.com/McFlanky/dreampic-ai/pkg/util"
	"github.com/McFlanky/dreampic-ai/view/auth"
	"github.com/nedpals/supabase-go"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if !util.IsValidEmail(credentials.Email) {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Email: "Please enter a valid email",
		}))
	}
	if reason, ok := util.IsValidPassword(credentials.Password); !ok {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			Password: reason,
		}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
