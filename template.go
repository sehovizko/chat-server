package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template // templ represents a single template.
}

// Serve the HTTP requests.
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

	log.Debug("templateHandler.ServeHTTP: Execute the template.")
	if err := t.templ.Execute(w, data); err != nil {
		log.Error("templateHandler.ServeHTTP: ", err)
	}
}
