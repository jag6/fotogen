package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jag6/fotogen/context"
	"github.com/jag6/fotogen/errors"
	"github.com/jag6/fotogen/models"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
		Search         Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Username string
	}
	data.Email = r.FormValue("email")
	data.Username = r.FormValue("username")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Username string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")
	if data.Email != "" && data.Username != "" && data.Password != "" {
		user, err := u.UserService.Create(data.Email, data.Username, data.Password)
		if err != nil {
			var modelError error
			switch errors.Is(err, modelError) {
			case modelError == models.ErrEmailTaken:
				err = errors.Public(err, "That email address is already associated with an account.")
			case modelError == models.ErrUsernameTaken:
				err = errors.Public(err, "That username is already taken.")
			}
			u.Templates.New.Execute(w, r, data, err)
			return
		}
		// if errors.Is(err, models.ErrEmailTaken) {
		// 	err = errors.Public(err, "That email address is already associated with an account.")
		// }
		// if errors.Is(err, models.ErrUsernameTaken) {
		// 	err = errors.Public(err, "That username is already taken.")
		// }
		session, err := u.SessionService.Create(user.ID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}
		setCookie(w, CookieSession, session.Token)
		http.Redirect(w, r, "/galleries", http.StatusFound)
	} else {
		http.Error(w, "Please fill out the form", http.StatusInternalServerError)
	}
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err, "User not found.")
		}
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "current user: %s\n", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/sign-in", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. For instance,
		// if a user doesn't exist with the email address.
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetUrl := "https://fotogenrfw.site/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// Don't render the token here! We need them to confirm they have access to their email in order to get the token.
	// Sharing it here would be a massive security hole.
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Sign the user in now that they have reset their password.
	// Any errors from this point onward should redirect to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/sign-in", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RedirectIfSignedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user != nil {
			http.Redirect(w, r, "/galleries", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) Search(w http.ResponseWriter, r *http.Request) {
	type User struct {
		ID       int
		Username string
	}
	var data struct {
		Users []User
	}
	q := r.URL.Query().Get("q")
	users, err := u.UserService.SearchByUsername(q)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, user := range users {
		data.Users = append(data.Users, User{
			ID:       user.ID,
			Username: user.Username,
		})
	}
	u.Templates.Search.Execute(w, r, data)
}
