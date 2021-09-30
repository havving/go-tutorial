package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	r := &router{make(map[string]map[string]HandlerFunc)}

	r.HandleFunc("GET", "/", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "Hello, World!")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	r.HandleFunc("GET", "/about", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "About")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	r.HandleFunc("GET", "/user/:id", logHandler(recoverHandler(func(c *Context) {
		if c.Params["id"] == "0" {
			fmt.Fprintf(c.ResponseWriter, "Get User %v\n", c.Params["id"])
		}
	})))

	r.HandleFunc("GET", "/user/:user_id/address/:address_id", func(c *Context) {
		fmt.Fprintf(c.ResponseWriter, "Get User %v's Address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	r.HandleFunc("POST", "/users", logHandler(recoverHandler(parseFormHandler(parseJsonBodyHandler(func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})))))

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		t := time.Now()
		fmt.Fprintf(c.ResponseWriter, "Create User %v's Address %v\n", c.Params["user_id"], c.Params["address_id"])
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	// 8080 포트로 웹 서버 구동
	http.ListenAndServe(":8080", r)
}
