package main

import (
	"fmt"
	"net/http"

	"github.com/gregpechiro/cook"
)

var count int

func main() {

	http.HandleFunc("/getflash", GetFlash)
	http.HandleFunc("/setflash", SetFlash)
	http.HandleFunc("/getall", GetAllCook)
	http.HandleFunc("/get", GetCook)
	http.HandleFunc("/set", SetCook)
	http.HandleFunc("/del", DelCook)

	http.ListenAndServe(":8080", nil)

}

func SetFlash(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("kind") {
	case "success":
		cook.SetSuccessRedirect(w, r, "/getflash", r.FormValue("msg"))
	case "error":
		cook.SetErrorRedirect(w, r, "/getflash", r.FormValue("msg"))
	default:
		cook.SetMsgRedirect(w, r, "/getflash", r.FormValue("msg"))
	}
	return
}

func GetFlash(w http.ResponseWriter, r *http.Request) {
	msgk, msgv := cook.GetFlash(w, r)

	fmt.Fprintf(w, "Message in cookie: %s, %s", msgk, msgv)
}

func GetAllCook(w http.ResponseWriter, r *http.Request) {
	c := cook.GetAllCookies(r)
	fmt.Fprintf(w, "Cookies:    %v", c)
}

func GetCook(w http.ResponseWriter, r *http.Request) {
	c := r.FormValue("cook")
	var val string
	if c == "" {
		val = ""
	} else {
		val = cook.GetCookie(r, r.FormValue("cook"))
	}
	fmt.Fprintf(w, "Cookie: %v", val)
}

func DelCook(w http.ResponseWriter, r *http.Request) {
	cook.DeleteCookie(w, r, r.FormValue("cook"))
	http.Redirect(w, r, "/getall", 303)
}

func SetCook(w http.ResponseWriter, r *http.Request) {
	count++
	key := fmt.Sprintf("KEY %d", count)
	val := fmt.Sprintf("VALUE %d", count)
	cook.PutCookie(w, key, val)
	fmt.Fprintf(w, "Set Cookie: %s, %s", key, val)
	return
}
