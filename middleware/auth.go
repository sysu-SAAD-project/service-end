package middleware

import (
	"strings"
	"net/http"

	"github.com/sysu-saad-project/service-end/controller"
	"github.com/sysu-saad-project/service-end/models/service"
)

// ValidUserMiddleWare check token and decide the permission
// 0->Basic permission
// 1->Activity uploader
// Timeout authorization will be set to 0
// Suppose that only manager can access the api
// User level and account name will be set
type ValidUserMiddleWare struct{}

func (v ValidUserMiddleWare) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	role := 0
	// Read authorization from header
	r.Header.Del("X-Role")
	r.Header.Del("X-Account")
	auth := r.Header.Get("Authorization")
	// Check user identity
	if len(auth) <= 0 {
		r.Header.Set("X-Role", "0")
	} else {
		ok, name := controller.CheckToken(auth)
		if ok != 2 {
			r.Header.Set("X-Role", "0")
		} else {
			// Check user account
			exist := service.IsUserExist(name)
			if !exist {
				r.Header.Set("X-Role", "0")
			} else {
				role = 1
				r.Header.Set("X-Role", "1")
				r.Header.Set("X-Account", name)
			}
		}
	}
	path := r.URL.Path
	if strings.HasPrefix(path, "/actApplys/") || path == "/actApplys" {
		if role == 0 {
			rw.WriteHeader(401)
			return
		}
	} else if path == "/discus" || strings.HasPrefix(path, "/discuss/") {
		if role == 0 {
			rw.WriteHeader(401)
			return
		}
	}
	next(rw, r)
}
