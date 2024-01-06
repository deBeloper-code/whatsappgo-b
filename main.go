package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// This type is responsible for loading,
// compiling and delivering our template.
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// * Templates allow us to blend generic text with specific text
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// With one Do only will be called once
	// Compiling a template is process by which the source template is INTERPRETED and PREPARED
	// for blending with various data.
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	// New room
	r := newRoom()
	//Handlers
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	//get the room going
	go r.run()
	// Starting server
	http.ListenAndServe(":8080", nil)
}
