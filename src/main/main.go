package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := &router{make(map[string]map[string]HandlerFunc)}

	r.HandleFunc("GET", "/", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "Hello, World!")
	})

	r.HandleFunc("GET", "/about", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "About")
	})

	r.HandleFunc("GET", "/user/:id", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "Get User %v\n", c.Params["id"])
	})

	r.HandleFunc("GET", "/user/:user_id/address/:address_id", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "Get User %v's Address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	r.HandleFunc("POST", "/users", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "Create User\n")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "Create User %v's Address\n", c.Params["user_id"])
	})

	// 8080 포트로 웹 서버 구동
	http.ListenAndServe(":8080", r)
}
