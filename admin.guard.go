package main

import (
	"net/http"
	"ro-backend/configuration"
	"strings"

	"github.com/golang-jwt/jwt"
)

func adminGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(bearerParts) < 2 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		jwtTokenStr := bearerParts[1]

		var jwtSecret = configuration.Config.Jwt.Secret

		claims := jwt.StandardClaims{}
		_, err := jwt.ParseWithClaims(jwtTokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if claims.Subject != "admin" {
			http.Error(w, "Only Admin", http.StatusForbidden)
			return
		}

		r.Header.Add("userId", claims.Id)

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}
