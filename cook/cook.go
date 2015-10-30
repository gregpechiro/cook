package cook

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var EXPYEAR = 3
var EXPMONTH = 0
var EXPDAY = 0

func Expires() time.Time {
	return time.Now().AddDate(EXPYEAR, EXPMONTH, EXPDAY)
}

type Cook struct {
	name, value string
}

func Get(r *http.Request, name string) string {
	cookie, err := GetCookie(r, name)
	if err != nil || cookie == nil {
		return ""
	}
	return BaseDec(cookie.Value)
}

func Put(w http.ResponseWriter, name, value string) {
	cookie := FreshCookie(name, value, Expires())
	http.SetCookie(w, &cookie)
}

func Delete(w http.ResponseWriter, r *http.Request, name string) {
	cookie, err := GetCookie(r, name)
	if err != nil {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(UrlEnc(name))
}

func PutCookie(w http.ResponseWriter, c *http.Cookie) {
	http.SetCookie(w, c)
}

func DeleteCookie(w http.ResponseWriter, c *http.Cookie) {
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func GetAll(r *http.Request) []Cook {
	return getAll(r.Cookies())
}

func getAll(cookies []*http.Cookie) []Cook {
	var cooks []Cook
	for _, c := range cookies {
		cooks = append(cooks, Cook{UrlDec(c.Name), BaseDec(c.Value)})
	}
	return cooks
}

func GetStartsWith(r *http.Request, pre string) []Cook {
	var cookies []*http.Cookie
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookies = append(cookies, cookie)
		}
	}
	return getAll(cookies)
}

func UpdateStartsWith(w http.ResponseWriter, r *http.Request, pre string, exp time.Time) {
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookie.Expires = exp
			cookie.Path = "/"
			cookie.HttpOnly = true
			http.SetCookie(w, cookie)
		}
	}
}

func DeleteStartsWith(w http.ResponseWriter, r *http.Request, pre string) {
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(UrlDec(cookie.Name), pre) {
			cookie.Expires = time.Now()
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
	}
}

func FreshCookie(name, value string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     UrlEnc(name),
		Value:    BaseEnc(value),
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	}
}

func BaseEnc(msg string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(msg))
}

func BaseDec(data string) string {
	msg, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return ""
	}
	return string(msg)
}

func UrlEnc(msg string) string {
	return url.QueryEscape(msg)
}

func UrlDec(data string) string {
	msg, err := url.QueryUnescape(data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return ""
	}
	return string(msg)
}
