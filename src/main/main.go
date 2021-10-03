package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

func main() {
	// 서버 생성
	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "Hello, World!")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
		c.RenderTemplate("/src/public/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *Context) {
		t := time.Now()
		fmt.Fprintln(c.ResponseWriter, "About")
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	s.HandleFunc("GET", "/user/:id", logHandler(recoverHandler(func(c *Context) {
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
		fmt.Fprintf(c.ResponseWriter, "Get User %v\n", c.Params["id"])
	})))

	s.HandleFunc("GET", "/user/:user_id/address/:address_id", func(c *Context) {
		u := User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJson(u)
		fmt.Fprintf(c.ResponseWriter, "Get User %v's Address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	s.HandleFunc("POST", "/users", logHandler(recoverHandler(parseFormHandler(parseJsonBodyHandler(func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})))))

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		t := time.Now()
		fmt.Fprintf(c.ResponseWriter, "Create User %v's Address %v\n", c.Params["user_id"], c.Params["address_id"])
		log.Printf("[%s] %q %v\n", c.Request.Method, c.Request.URL.String(), time.Now().Sub(t))
	})

	s.HandleFunc("GET", "/login", func(c *Context) {
		// "login.html" 렌더링
		c.RenderTemplate("/src/public/login.html",
			map[string]interface{}{"message": "로그인이 필요합니다."})
	})

	s.HandleFunc("POST", "/login", func(c *Context) {
		// 로그인 정보를 확인하여 쿠키에 인증 토큰 값 기록
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		// id와 password가 맞지 않으면 다시 "/login" 페이지 렌더링
		c.RenderTemplate("/src/public/login.html",
			map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다."})
	})

	s.Use(AuthHandler)

	// 웹 서버 구동
	s.Run(":8080")
}

/** 인증된 웹 요청만 허용하는 미들웨어 **/
const VerifyMessage = "verified"

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "/src/public/index.html"}
	return func(c *Context) {
		// url prefix가 "/login", "/src/public/index.html"이면 auth를 체크하지 않음
		for _, s := range ignore {
			if strings.HasSuffix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// "X_AUTH" 쿠키 값이 없으면 "/login"으로 이동
			c.Redirect("/login")
			return
		} else if err != nil {
			// 에러 처리
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			// 쿠키 값으로 인증이 확인되면 다음 핸들러로 넘어감
			next(c)
			return
		}

		// "/login"으로 이동
		c.Redirect("/login")
	}
}

/** 인증 토큰 확인 **/
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

/** 로그인 처리 **/
func CheckLogin(username, password string) bool {
	const (
		USERNAME = "tester"
		PASSWORD = "1234"
	)

	return username == USERNAME && password == PASSWORD
}

/** 인증 토큰 생성 **/
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return " "
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)

	return hex.EncodeToString(mac.Sum(nil))
}
