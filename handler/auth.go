package handler

import (
	"fmt"
	"net/http"

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
	// call supabase

	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
		InvalidCredentials: "The credentials you have entered are invalid",
	}))

	fmt.Println(credentials)
	return nil
}
