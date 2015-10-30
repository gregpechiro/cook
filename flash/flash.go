package flash

import (
	"net/http"
	"strings"

	"github.com/gregpechiro/cookieManager/cook"
)

func SetFlash(w http.ResponseWriter, kind, msg string) {
	cook.Put(w, "flash", kind+":"+msg)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, string) {
	msg := strings.Split(cook.Get(r, "flash"), ":")
	cook.Delete(w, r, "flash")
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
