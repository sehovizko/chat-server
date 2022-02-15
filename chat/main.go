package main

import (
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	log "github.com/sirupsen/logrus"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template  // templは1つのテンプレートを表す。
}

// HTTPリクエストを処理する。
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    log.Info("templateHandler.ServeHTTP: HTTP接続を開始します。")
    if err := t.templ.Execute(w, r); err != nil {
        log.Error("templateHandler.ServeHTTP: ", err)
    }
}

func main() {
	log.Info("main: 準備開始します。")
	r := newRoom()
    http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	log.Info("main: 準備完了しました。" )

	// チャットルームを開始する。
    go r.run()

	var addr = flag.String("addr", ":8080", " アプリケーションのアドレス")
	flag.Parse()  // フラグ解釈する。

	// Webサーバーを開始する。
	log.Info("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Error("ListenAndServe: ", err)
	}
}
