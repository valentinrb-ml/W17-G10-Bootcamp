package middleware

import (
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/apperrors"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/api/response"
	"net/http"
)

func AllowMethods(methods ...string) func(http.Handler) http.Handler {
	methodMap := make(map[string]bool, len(methods))
	for _, method := range methods {
		methodMap[method] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !methodMap[r.Method] {
				response.Error(w, apperrors.NewAppError(apperrors.CodeMethodNotAllowed, "method not allowed"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Helpers for common cases
var (
	OnlyGET    = AllowMethods("GET")
	OnlyPOST   = AllowMethods("POST")
	OnlyPATCH  = AllowMethods("PATCH")
	OnlyDELETE = AllowMethods("DELETE")
)
