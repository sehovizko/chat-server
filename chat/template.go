package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template // templは1つのテンプレートを表す。
}

// HTTPリクエストを処理する。
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	log.Info("templateHandler.ServeHTTP: テンプレートの処理を行います。")
	if err := t.templ.Execute(w, r); err != nil {
		log.Error("templateHandler.ServeHTTP: ", err)
	}
}
