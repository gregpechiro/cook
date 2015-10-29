package cook

import (
	"encoding/base32"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var EXPIRES = time.Now().AddDate(3, 0, 0)

type Cook struct {
	name, value string
}

func GetAllCookies(r *http.Request) []Cook {
	var cooks []Cook
	for _, c := range r.Cookies() {
		cooks = append(cooks, Cook{urlDec(c.Name), baseDec(c.Value)})
	}
	return cooks
}

func GetCookie(r *http.Request, name string) string {
	cookie, err := getCookie(r, name)
	if err != nil {
		return ""
	}
	return baseDec(cookie.Value)
}

func PutCookie(w http.ResponseWriter, name, value string) {
	cookie := freshCookie(name, value, EXPIRES)
	http.SetCookie(w, &cookie)
}

func DeleteCookie(w http.ResponseWriter, r *http.Request, name string) {
	cookie, err := getCookie(r, name)
	if err != nil {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

func SetFlash(w http.ResponseWriter, kind, msg string) {
	cookie := msgCookie(fmt.Sprintf("%s:%s", kind, msg))
	http.SetCookie(w, &cookie)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	cookie, err := getCookie(r, "msg")
	if err != nil {
		return "", ""
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	msg := strings.Split(baseDec(cookie.Value), ":")
	if len(msg) != 2 {
		return "", ""
	}
	return msg[0], msg[1]
}

func SetSuccessRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alertSuccess", msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetErrorRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alertError", msg)
	http.Redirect(w, r, url, 303)
	return
}

func SetMsgRedirect(w http.ResponseWriter, r *http.Request, url, msg string) {
	SetFlash(w, "alert", msg)
	http.Redirect(w, r, url, 303)
	return
}

func msgCookie(msg string) http.Cookie {
	return freshCookie("msg", msg, EXPIRES) // 3 years in the future
}

func freshCookie(name, value string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     urlEnc(name),
		Value:    baseEnc(value),
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
	}
}

func getCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(urlEnc(name))
}

func baseEnc(msg string) string {
	return base32.StdEncoding.EncodeToString([]byte(msg))
}

func baseDec(data string) string {
	msg, err := base32.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return ""
	}
	return string(msg)
}

func urlEnc(msg string) string {
	return url.QueryEscape(msg)
}

func urlDec(data string) string {
	msg, err := url.QueryUnescape(data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return ""
	}
	return string(msg)
}
