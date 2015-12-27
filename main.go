package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once  sync.Once
	file  string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.file)))
	})

	t.templ.Execute(w, r)
}

func main() {
	log.SetPrefix("[chat] - ")
	var port = flag.String("port", "8080", "The port of the application.")
	flag.Parse()

	r := newRoom()
	http.Handle("/login", &templateHandler{file: "login.html"})
	http.Handle("/", MustAuth(&templateHandler{file: "chat.html"}))
	http.Handle("/room", MustAuth(r))
	http.HandleFunc("/auth/", loginHandler)
	//get the room going
	go r.run()
	log.Println("start server in  port:", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
