package middleware

import (
	"net/http"
	"strings"

	"github.com/Luiggy102/go-rest-ws/models"
	"github.com/Luiggy102/go-rest-ws/server"
	"github.com/golang-jwt/jwt"
)

func CheckAuth(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				// start middleware
				currentPath := r.URL.Path

				// check if the path need an auth
				noAuthNeededPath := []string{"login", "signup"}
				for _, p := range noAuthNeededPath {
					// if the current path don't need an auth
					if !strings.Contains(currentPath, p) {
						// run the next handler
						next.ServeHTTP(w, r)
						return
					}
				}
				// current path need an authorization token
				// get the token
				token := r.Header.Get("Authorization")
				tokenStr := strings.TrimSpace(token)

				// check the token

				// jwt needs a func to get the token secret key
				var keyfunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
					// the used to sign the token
					return []byte(s.Config().JWTSecret), nil
				}

				// parse the token
				// if error the token is not valid
				_, err := jwt.ParseWithClaims(tokenStr, models.AppClaims{}, keyfunc)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				// token is valid
				next.ServeHTTP(w, r)
			},
		)
	}
}
