package main

import (
	"net/http"
	"strings"
)

/** 라우터 타입 정의 **/
type router struct {
	// key: http method
	// value: url 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]HandlerFunc
}

/** 라우터에 핸들러를 등록하기 위한 메서드 **/
func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	// http method로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵 생성
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}
	// http method로 등록된 맵에 url 패턴과 핸들러 함수 등록
	m[pattern] = h
}

/** http.Handler 인터페이스의 ServeHTTP 메서드 정의 **/
type Handler interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}

/** http 메서드와 url 경로를 분석하여 맞는 핸들러를 찾아 동작 시키는 함수**/
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// http 메서드에 맞는 모든 handers를 반복하여 요청 url에 해당하는 handler를 찾음
	for pattern, handler := range r.handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {
			// Context 생성
			c := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: w,
				Request:        req,
			}
			for k, v := range params {
				c.Params[k] = v
			}

			// 요청 url에 해당하는 handler 수행
			handler(&c)
			return
		}
	}
	// 요청 url에 해당하는 handler를 찾지 못하면 NotFound 에러 반환
	http.NotFound(w, req)
	return
}

/** 라우터에 등록된 동적 url 패턴과 실제 url 경로가 일치하는지 확인하는 함수 **/
func match(pattern, path string) (bool, map[string]string) {
	// pattern과 path가 정확히 일치하면 true 반환
	if pattern == path {
		return true, nil
	}

	// pattern과 path를 "/" 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// pattern과 path를 "/"로 구분한 후, 부분 문자열 집합의 개수가 다르면 false 반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	// pattern에 일치하는 url 매개변수를 담기 위한 params 맵 생성
	params := make(map[string]string)

	// "/"로 구분된 pattern/path의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// pattern과 path의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// pattern이 ':'로 시작하면 params에 url params를 다음 후 다음 루프
			params[patterns[i][1:]] = paths[i]
		default:
			// 일치하는 경우가 없으면 false 반환
			return false, nil
		}
	}
	// true와 params 반환
	return true, params
}

/** 핸들러 함수 반환 **/
func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		// http 메서드에 맞는 모든 handlers를 반복하며 요청 url에 해당하는 handler를 찾음
		for pattern, handler := range r.handlers[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}
				// 요청 url에 해당하는 handler 수행
				handler(c)
				return
			}
		}
		// 요청 url에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}
