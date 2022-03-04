package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Global variables
var (
	log  = logrus.New()
	addr = flag.String("addr", ":8080", " アプリケーションのアドレス")
)

func init() {
	// フラグ解釈する。
	flag.Parse()

	// Gothのセットアップ
	goth.UseProviders(
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), fmt.Sprintf("http://localhost%s/auth/gplus/callback", *addr)),
	)

	// ログレベルを決定する。
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	log.Debug("main: ルーティングを開始します。")

	router := pat.New()
	router.Get("/auth/{provider}/callback", loginCallbackHandler)
	router.Get("/auth/{provider}", loginHandler)
	router.Get("/logout", logoutHandler)

	router.Add("GET", "/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	router.Add("GET", "/login", &templateHandler{filename: "login.html"})

	r := newRoom()
	router.Add("GET", "/room", r)

	log.Debug("main: ルーティングを終了しました。")

	// チャットルームを開始する。
	go r.run()

	// Webサーバーを開始する。
	log.Info("Webサーバーを開始します。ポート: ", *addr)
	log.Fatal(http.ListenAndServe(*addr, router))
}
