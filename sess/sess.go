package sess

import (
	"net/http"
	"strings"
	"time"

	"github.com/gregpechiro/cookieManager/cook"
)

var SESSDUR = time.Minute * 15

func toMap(data string) map[string]string {
	m := make(map[string]string)
	for _, kv := range strings.Split(data, "\n") {
		set := strings.Split(kv, "^")
		if len(set) == 2 {
			m[set[0]] = set[1]
		}
	}
	return m
}

func toString(m map[string]string) string {
	var ss []string
	for k, v := range m {
		ss = append(ss, k+"^"+v)
	}
	return strings.Join(ss, "\n")
}

func SessDur() time.Time {
	return time.Now().Add(SESSDUR)
}

func Login(w http.ResponseWriter, r *http.Request) {
	m := toMap(cook.Get(r, "SESS-D"))
	if m == nil {
		m = make(map[string]string)
	}
	m["AUTH"] = "true"
	cookie := cook.FreshCookie("SESS-D", toString(m), SessDur())
	cook.PutCookie(w, &cookie)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cook.Delete(w, r, "SESS-D")
}

func Put(w http.ResponseWriter, r *http.Request, name, val string) {
	m := toMap(cook.Get(r, "SESS-D"))
	m[name] = val
	cookie := cook.FreshCookie("SESS-D", toString(m), SessDur())
	cook.PutCookie(w, &cookie)
}

func GetAll(r *http.Request) map[string]string {
	return toMap(cook.Get(r, "SESS-D"))
}

func Authorized(w http.ResponseWriter, r *http.Request) bool {
	m := toMap(cook.Get(r, "SESS-D"))
	if m == nil || m["AUTH"] != "true" {
		return false
	}
	cookie := cook.FreshCookie("SESS-D", toString(m), SessDur())
	cook.PutCookie(w, &cookie)
	return true
}
