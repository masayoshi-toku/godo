package main

import (
	"fmt"
	"net/http"
)

type WelcomeHandler struct{}

func (h *WelcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to GODO!")
}

func main() {
	welcomeHandler := WelcomeHandler{}

	server := http.Server{
		Addr: ":8080",
	}

	http.Handle("/", &welcomeHandler)
	server.ListenAndServe()
}
