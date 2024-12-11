package cookies

import (
	"net/http"
	"time"
)

func GetSessionCookie(name string, r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil
	}
	return cookie
}

func SetSessionCookie(name string, w http.ResponseWriter, token string, exptime time.Time) {
	cookie := http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		Expires:  exptime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func ExpireSessionCookie(name string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}
