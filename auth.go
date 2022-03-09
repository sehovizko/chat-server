package main

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// If it was unauthentication or no cookie existed.
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		log.Info("authHandler.serveHTTP: Unauthenticated.")
	} else if err != nil {
		// Other errors occurred.
		log.Error(err)
		panic(err.Error())
	} else {
		// If the authentication succeeds, call the wrapped handler.
		h.next.ServeHTTP(w, r)
		log.Info("Authentication succeded.")
	}
}

// MustAuth forces user to be authenticated.
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("loginHandler: login handler was invoked.")

	// Call the authentication handler.
	gothic.BeginAuthHandler(w, r)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("loginCallbackHandler: Start login handler")

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil { // Call the authentication handler.
		log.Warning(fmt.Fprintln(w, err))
		return
	}

	authCookieValue := objx.New(map[string]interface{}{
		"name":       user.Name,
		"avatar_url": user.AvatarURL,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	// Go to chat.
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)

	log.Debug("loginCallbackHandler: End login handler.")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("logoutHandler: logout handler was invoked.")

	// Delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Setting MaxAge to -1, which indicates that it should be deleted immediately by the browser.
	})

	gothic.Logout(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
