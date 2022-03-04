package main

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未承認だった場合
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		log.Info("authHandler.serveHTTP: 未認証です。")
	} else if err != nil {
		// 別の何らかのエラーが発生
		log.Error(err)
		panic(err.Error())
	} else {
		// 認証に成功した場合、ラップされたハンドラを呼び出す。
		h.next.ServeHTTP(w, r)
		log.Info("認証に成功しました。")
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("loginHandler: ログインハンドラが呼び出されました。")

	// 承認ハンドラを呼び出します。
	gothic.BeginAuthHandler(w, r)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("loginCallbackHandler: ログインコールバック開始します。")

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil { // 何らかの理由でユーザー認証が完了しなかった。
		log.Warning(fmt.Fprintln(w, err))
		return
	}

	authCookieValue := objx.New(map[string]interface{}{
		"name": user.Name,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	// チャット画面へ
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)

	log.Debug("loginCallbackHandler: ログインコールバック終了します。")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("logoutHandler: ログアウトハンドラが呼び出されました。")

	http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Value: "",
		Path: "/",
		MaxAge: -1,
	})

	gothic.Logout(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
