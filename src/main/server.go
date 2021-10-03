package main

import "net/http"

/** Server 타입 정의 **/
type Server struct {
	*router
	middlewares  []Middleware
	startHandler HandlerFunc
}

/** 생성자 함수로 Server 생성 **/
func NewServer() *Server {
	r := &router{make(map[string]map[string]HandlerFunc)}
	s := &Server{router: r}
	s.middlewares = []Middleware{
		logHandler,
		recoverHandler,
		staticHandler,
		parseFormHandler,
		parseJsonBodyHandler}
	return s
}

/** 웹 서버 구동 메서드 **/
func (s *Server) Run(addr string) {
	// startHandler를 라우터 핸들러 함수로 지정
	s.startHandler = s.router.handler()

	// 등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	// 웹 서버 시작
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}

/** Server가 웹 요청을 받아 처리하는 메서드 **/
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Context 생성
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	s.startHandler(c)
}

/** 새 미들웨어 추가 **/
func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}
