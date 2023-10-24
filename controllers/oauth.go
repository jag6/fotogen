package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}
	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)
	url := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
		// oauth2.SetAuthURLParam("token_access_type", "offline"),
		// oauth2.AccessTypeOffline,
	)
	http.Redirect(w, r, url, http.StatusFound)
}

func (oa OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)
	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}
	state := r.FormValue("state")
	cookieState, err := readCookie(r, "oauth_state")
	if err != nil || cookieState != state {
		if err != nil {
			fmt.Println(err)
		}
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	deleteCookie(w, "oauth_state")
	code := r.FormValue("code")
	token, err := config.Exchange(
		r.Context(),
		code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI(r, provider)),
	)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	// Persist the user's oauth token so we can use it in the future...
	// For now, we can just print it out.

	// w.Header().Set("Content-type", "application/json")
	// enc := json.NewEncoder(w)
	// enc.SetIndent("", " ")
	// enc.Encode(token)

	client := config.Client(r.Context(), token)
	res, err := client.Post(
		"https://api.dropboxapi.com/2/files/list_folder", "application/json",
		strings.NewReader(`{"path": ""}`),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	unpretty, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var pretty bytes.Buffer
	err = json.Indent(&pretty, unpretty, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	pretty.WriteTo(w)
}

func redirectURI(r *http.Request, provider string) string {
	if r.Host == "localhost:3000" {
		return fmt.Sprintf("http://localhost:3000/oauth/%s/callback", provider)
	}
	return fmt.Sprintf("https://fotogenrfw.online/oauth/%s/callback", provider)
}
