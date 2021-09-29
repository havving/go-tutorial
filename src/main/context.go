package main

import "net/http"

/** 컨텍스트 타입 정의 **/
type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

/** 핸들러 타입 정의 **/
type HandlerFunc func(*Context)
