package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gregpechiro/cookieManager/cook"
)

func SetCSRF(w http.ResponseWriter) string {
	csrf := Random(32)
	cook.Put(w, "C", csrf)
	return csrf
}

func GetCSRF(r *http.Request) string {
	return cook.Get(r, "C")
}

func ValidCSRF(r *http.Request) bool {
	c := GetCSRF(r)
	tok := r.FormValue("_csrf")
	return (tok == c && c != "")
}

func Random(n int) string {
	e := make([]byte, n)
	rand.Read(e)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(e)))
	base64.URLEncoding.Encode(b, e)
	return string(b)[:n]
}
