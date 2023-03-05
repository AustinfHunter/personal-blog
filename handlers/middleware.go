package handlers

import (
	"net/http"
)

func CorsMiddleWare(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
	}
}

func authorizationMiddleware(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	pt := parseToken(auth)

	return pt.Valid
}
