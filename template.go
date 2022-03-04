package main

import (
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	} else {
		log.Error("templateHandler.ServeHTTP: ", err)
	}

	log.Debug("templateHandler.ServeHTTP: テンプレートの処理を行います。")
	if err := t.templ.Execute(w, data); err != nil {
		log.Error("templateHandler.ServeHTTP: ", err)
	}
}
