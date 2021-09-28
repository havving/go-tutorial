package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := &router{make(map[string]map[string]http.HandlerFunc)}

	r.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	r.HandleFunc("GET", "/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "About")
	})

	r.HandleFunc("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Get User")
	})

	r.HandleFunc("GET", "/user/:user_id/address/:address_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Get User's Address")
	})

	r.HandleFunc("POST", "/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create User")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create User's Address")
	})

	// 8080 포트로 웹 서버 구동
	http.ListenAndServe(":8080", r)
}
