package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

// This type is responsible for loading,
// compiling and delivering our template.
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// * Templates allow us to blend generic text with specific text
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// With one Do only will be called once
	// Compiling a template is process by which the source template is INTERPRETED and PREPARED
	// for blending with various data.
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("..", "templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New("396031989690-38rm4p9ftrvi19nb7fejkds5fkkb4qf0.apps.googleusercontent.com", "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",
			"http://localhost:8080/auth/callback/google"),
	)
	// New room
	r := newRoom()
	//Handlers
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//get the room going
	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
