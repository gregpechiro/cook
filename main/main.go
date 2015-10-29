package main

import (
	"fmt"
	"net/http"

	"github.com/gregpechiro/cook"
)

var count int

func main() {

	http.HandleFunc("/get", GetFlash)
	http.HandleFunc("/set", SetFlash)

	http.HandleFunc("/getcook", GetCook)
	http.HandleFunc("/setcook", SetCook)
	http.HandleFunc("/delcook", DelCook)

	http.ListenAndServe(":8080", nil)

}

func SetFlash(w http.ResponseWriter, r *http.Request) {
	//cook.SetFlash(w, "alertSuccess", r.FormValue("msg"))
	//http.Redirect(w, r, "/get", 303)
	//return
	switch r.FormValue("kind") {
	case "success":
		cook.SetSuccessRedirect(w, r, "/get", r.FormValue("msg"))
	case "error":
		cook.SetErrorRedirect(w, r, "/get", r.FormValue("msg"))
	default:
		cook.SetMsgRedirect(w, r, "/get", r.FormValue("msg"))
	}
	return
}

func GetFlash(w http.ResponseWriter, r *http.Request) {
	msgk, msgv := cook.GetFlash(w, r)

	fmt.Fprintf(w, "Message in cookie: %s, %s", msgk, msgv)
}

func GetCook(w http.ResponseWriter, r *http.Request) {
	c := cook.GetAllCookies(r)
	fmt.Fprintf(w, "Cookies:    %v", c)
}

func DelCook(w http.ResponseWriter, r *http.Request) {
	cook.DeleteCookie(w, r, r.FormValue("cook"))
	http.Redirect(w, r, "/getcook", 303)
}

func SetCook(w http.ResponseWriter, r *http.Request) {
	count++
	key := fmt.Sprintf("KEY_%d", count)
	val := fmt.Sprintf("VALUE_%d", count)
	cook.PutCookie(w, key, val)
	fmt.Fprintf(w, "Set Cookie: %s, %s", key, val)
	return
}
