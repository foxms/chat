package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 3 {
		log.Println("Auth URL length error: ", r.URL.Path)
		return
	}
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	case "direct":
		user := r.PostFormValue("user")
		if len(user) < 1 || user == "" {
			log.Println("Get user name failed., try login again!")
			w.Header().Set("Location", "/login")
		} else {
			log.Println("User [", user, "] login")
			usr := base64.StdEncoding.EncodeToString([]byte(user))
			http.SetCookie(w, &http.Cookie{
				Name:  "auth",
				Value: usr,
				Path:  "/"})

			w.Header().Set("Location", "/chat")
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		log.Println(w, "Auth action ", action, " not supported")
	}
}
