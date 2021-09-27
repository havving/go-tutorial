package main

import "net/http"

/** 라우터 타입 정의 **/
type router struct {
	// key: http method
	// value: url 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]http.HandlerFunc
}

/** 라우터에 핸들러를 등록하기 위한 메서드 **/
func (r *router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	// http method로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵 생성
		m = make(map[string]http.HandlerFunc)
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
	if m, ok := r.handlers[req.Method]; ok {
		if h, ok := m[req.URL.Path]; ok {
			// 요청 url에 해당하는 핸들러 수행
			h(w, req)
			return
		}
	}
	// 요청에 일치하는 핸들러가 등록되어 있지 않으면 NotFound 에러 반환
	http.NotFound(w, req)
}
